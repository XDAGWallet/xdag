package wallet

import (
	"encoding/hex"
	"fmt"
	"goXdagWallet/xdago/common"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/magiconair/properties/assert"

	"goXdagWallet/xdago/secp256k1"
)

const (
	PRIVATE_KEY_STRING = "a392604efc2fad9c0b3da43b5f698a2e3f270f170d859912be0d54742275c5f6"
	PUBLIC_KEY_STRING  = "0x506bc1dc099358e5137292f4efdd57e400f29ba5132aa5d12b18dac1c1f6aab" +
		"a645c0b7b58158babbfa6c6cd5a48aa7340a8749176b120e8516216787a13dc76"
	PUBLIC_KEY_COMPRESS_STRING = "02506bc1dc099358e5137292f4efdd57e400f29ba5132aa5d12b18dac1c1f6aaba"
	ADDRESS                    = "b731bf10ed204f4ebc3d32ac88b7aa61b993fd59"
	PASSWORD                   = "Insecure Pa55w0rd"
	MNEMONIC                   = "scatter major grant return flee easy female jungle" +
		" vivid movie bicycle absent weather inspire carry"
)

func setup() (string, *Wallet) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	pwd := "password"
	wallet := NewWallet(path.Join(common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	wallet.UnlockWallet(pwd)
	keyBytes, _ := hex.DecodeString(PRIVATE_KEY_STRING)
	privKey := secp256k1.PrivKeyFromBytes(keyBytes)
	wallet.AddAccount(privKey)
	wallet.Flush()
	wallet.LockWallet()
	return pwd, &wallet
}

func TestGetPassword(t *testing.T) {
	p, w := setup()
	defer tearDown(w) //

	w.UnlockWallet(p)
	assert.Equal(t, w.GetPassword(), "password")
}

func TestUnlock(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	assert.Equal(t, w.IsUnLocked(), false)

	w.UnlockWallet(p)
	assert.Equal(t, w.IsUnLocked(), true)

	assert.Equal(t, len(w.GetAccounts()), 1)
}

func TestLock(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)
	assert.Equal(t, w.mnemonicPhrase, "")
	w.LockWallet()
	assert.Equal(t, w.IsUnLocked(), false)
}

func TestAddAccounts(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)
	w.SetAccounts([]*secp256k1.PrivateKey{})
	key1, _ := secp256k1.GeneratePrivateKey()
	key2, _ := secp256k1.GeneratePrivateKey()
	w.SetAccounts([]*secp256k1.PrivateKey{key1, key2})
	accounts := w.GetAccounts()
	assert.Equal(t, accounts[0].Key.String(), key1.Key.String())
	assert.Equal(t, accounts[1].Key.String(), key2.Key.String())

}

func TestFlush(t *testing.T) {
	p, w := setup()
	defer tearDown(w)
	info, _ := os.Stat(w.GetFile())
	size := info.Size()
	fmt.Println(size)
	w.UnlockWallet(p)
	w.SetAccounts([]*secp256k1.PrivateKey{})

	w.Flush()

	info2, _ := os.Stat(w.GetFile())
	size2 := info2.Size()

	if size2 >= size {
		panic("wallet file error")
	}
}

func TestFChangePassword(t *testing.T) {
	pwd2 := "passw0rd2"
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)
	w.ChangePassword(pwd2)
	w.Flush()
	w.LockWallet()

	assert.Equal(t, w.UnlockWallet(p), false)
	assert.Equal(t, w.UnlockWallet(pwd2), true)
}

func TestAccountRandom(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)
	oldSize := len(w.GetAccounts())
	w.AddAccountRandom()
	assert.Equal(t, len(w.GetAccounts()), oldSize+1)
}

func TestRemoveAccount(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)
	oldSize := len(w.GetAccounts())

	key, _ := secp256k1.GeneratePrivateKey()
	w.AddAccount(key)
	assert.Equal(t, len(w.GetAccounts()), oldSize+1)

	w.RemoveAccountByKey(key)
	assert.Equal(t, len(w.GetAccounts()), oldSize)

	w.AddAccount(key)
	assert.Equal(t, len(w.GetAccounts()), oldSize+1)

	w.RemoveAccountByKey(key)
	assert.Equal(t, len(w.GetAccounts()), oldSize)
}

func TestAddAccountWithNextHdKey(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)

	w.InitializeHdWallet(MNEMONIC)
	oldSize := len(w.GetAccounts())

	for i := 0; i < 5; i++ {
		w.AddAccountWithNextHdKey()
	}
	//log.Debug("account[4]", log.Ctx{"key": w.GetAccounts()[4].Serialize()})
	assert.Equal(t, len(w.GetAccounts()), oldSize+5)
}

func TestOnlyHdKey(t *testing.T) {
	p, w := setup()
	//defer tearDown(w)

	w.UnlockWallet(p)

	w.InitializeHdWallet(MNEMONIC)

	assert.Equal(t, len(w.GetAccounts()), 1)

	key := w.GetDefKey()
	w.RemoveAccountByKey(key)
	assert.Equal(t, len(w.GetAccounts()), 0)

	w.AddAccountWithNextHdKey()
	assert.Equal(t, len(w.GetAccounts()), 1)

	//w.InitializeHdWallet("")  // remove mnemonic and test wallet file size is 125 bytes

	w.Flush()
	w.LockWallet()

}

func TestHDKeyRecover(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)

	w.InitializeHdWallet(MNEMONIC)

	var keys1 [5]string
	var keys2 [5]string

	for i := 0; i < 5; i++ {
		key := w.AddAccountWithNextHdKey()
		keys1[i] = key.Key.String()
	}

	w2 := NewWallet(path.Join(common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	w2.UnlockWallet(p + p)
	w2.InitializeHdWallet(MNEMONIC)
	for i := 0; i < 5; i++ {
		key := w2.AddAccountWithNextHdKey()
		keys2[i] = key.Key.String()
	}
	assert.Equal(t, keys2, keys1)
}

func tearDown(w *Wallet) {
	w.Delete()
}

func TestNewMnemonic(t *testing.T) {
	mnemonic := NewMnemonic(128)
	fmt.Println(mnemonic)

}

func TestNewBipWallet(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	w := NewWallet(path.Join(common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	w.UnlockWallet("ljr20040224")
	mnemonic := NewMnemonic(128)
	w.InitializeHdWallet(mnemonic)
	w.AddAccountWithNextHdKey()
	fmt.Println(w.mnemonicPhrase)
	res := w.Flush()

	assert.Equal(t, res, true)
	assert.Equal(t, len(w.GetAccounts()), 1)

	w2 := NewWallet(path.Join(common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	w2.UnlockWallet("ljr20040224")
	assert.Equal(t, len(w2.GetAccounts()), 1)
	assert.Equal(t, w.mnemonicPhrase, mnemonic)
	//w.Delete()
}

func TestExportDefKey(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)

	err := w.ExportDefKey("priv.key")

	assert.Equal(t, err, nil)
}

func TestExportMnemonic(t *testing.T) {
	p, w := setup()
	defer tearDown(w)

	w.UnlockWallet(p)

	w.InitializeHdWallet(MNEMONIC)

	err := w.ExportMnemonic("mnemonic.txt")

	assert.Equal(t, err, nil)
}

func TestImportFromKey(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	pwd := "password"
	w, err := ImportWalletFromDefKey("priv.key", ".", pwd)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(w.GetAccounts()), 1)

	w2 := NewWallet(path.Join(common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	w2.UnlockWallet("password")
	assert.Equal(t, len(w2.GetAccounts()), 1)

	assert.Equal(t, w2.GetDefKey().Key.String(), PRIVATE_KEY_STRING)
	w.Delete()

}

func TestImportFromMnemonic(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	pwd := "password"
	w, err := ImportWalletFromMnemonicFile("mnemonic.txt", ".", pwd)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(w.GetAccounts()), 1)

	w2 := NewWallet(path.Join(common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	w2.UnlockWallet("password")
	assert.Equal(t, len(w2.GetAccounts()), 1)
	assert.Equal(t, w.mnemonicPhrase, MNEMONIC)
	w.Delete()

}
