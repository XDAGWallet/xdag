package server

import (
	"fmt"
	"goXdagWallet/cli"
	"goXdagWallet/components"
	"goXdagWallet/fileutils"
	"goXdagWallet/xdago/common"
	bip "goXdagWallet/xdago/wallet"
	"os"
	"path"
)

func ImportServWallet() error {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)

	filePath := cli.InputFilePath()
	cli.WalletAccount.Password = cli.ShowPassword()
	var err error

	pathDest := path.Join(pwd, common.BIP32_WALLET_FOLDER)
	if err = os.RemoveAll(pathDest); err != nil {
		fmt.Println("Clear dir failed", err)
		return err
	}
	if err = fileutils.MkdirAll(pathDest); err != nil {
		fmt.Println("Make dir failed", err)
		return err
	}
	components.BipWallet, err = bip.ImportWalletFromMnemonicFile(filePath, pwd, cli.WalletAccount.Password)
	if err != nil {
		fmt.Println("Import mnemonic failed", err)
		return err
	} else {
		components.BipWallet.LockWallet()
		components.PwdStr = ""
		components.BipAddress = ""
		return nil
	}
}
