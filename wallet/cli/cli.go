package cli

import (
	"fmt"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/fileutils"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xdago/cryptography"
	bip "goXdagWallet/xdago/wallet"
	"goXdagWallet/xlog"
	"os"
	"path"

	"github.com/manifoldco/promptui"
)

var WalletAccount components.WalletState

func NewCli(accountStatus int) {
	fmt.Println("Starting xdag wallet, version ", config.GetConfig().Version)
	xlog.Info("Starting xdag wallet, version ", config.GetConfig().Version)
	WalletAccount.HasAccount = accountStatus >= 0
	WalletAccount.WalletType = accountStatus
	if WalletAccount.HasAccount {
		connectWallet()
	} else {
		registerWallet()
	}
	OpenAndRunWallet()
}

func connectWallet() {
	fmt.Println("Connecting wallet....")
	if WalletAccount.WalletType == components.HAS_BOTH {
		selectWallet()
	} else {
		WalletAccount.Password = ShowPassword()
		copy(components.Password[:], WalletAccount.Password)
	}
}
func selectWallet() {
	prompt := promptui.Select{
		Label: "Select wallet type",
		Items: []string{"Mnemonic", "Non Mnemonic"},
	}

	_, result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input selection failed %v\n", err)
		_, result, err = prompt.Run()
	}

	if result == "Mnemonic" {
		WalletAccount.WalletType = components.HAS_ONLY_BIP
	} else {
		WalletAccount.WalletType = components.HAS_ONLY_XDAG
	}
	WalletAccount.Password = ShowPassword()
	copy(components.Password[:], WalletAccount.Password)
}
func registerWallet() {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	fmt.Println("Registering wallet....")
	prompt := promptui.Select{
		Label: "Select Wallet Type",
		Items: []string{"Create Mnemonic", "Import Mnemonic", "Import Non Mnemonic"},
	}

	_, result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		_, result, err = prompt.Run()
	}

	if result == "Import Non Mnemonic" {
		WalletAccount.WalletType = components.HAS_ONLY_XDAG
		folder := inputFolderPath()
		if components.CopyOldWallet(folder) != nil {
			fmt.Println("Import wallet data failed")
		} else {
			WalletAccount.Password = ShowPassword()
			copy(components.Password[:], WalletAccount.Password)
			res := components.ConnectXdagWallet()
			if res == 0 {
				OpenAndRunWallet()
			} else {
				fmt.Println("Password incorrect")
			}
		}
	} else {
		WalletAccount.WalletType = components.HAS_ONLY_BIP
		if result == "Create Mnemonic" {
			var psDone bool
			var passwd string
			for !psDone {
				passwd = ShowPassword()
				psDone = reshowPassword(passwd)
			}

			w, res := components.NewBipWallet(passwd, 128)
			if res {
				components.BipWallet = w
				b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
				components.BipAddress = base58.ChkEnc(b[:])
				components.PwdStr = passwd
				WalletAccount.Password = passwd
				OpenAndRunWallet()
			}

		} else { // Import Mnemonic
			filePath := InputFilePath()
			WalletAccount.Password = ShowPassword()

			pathDest := path.Join(pwd, common.BIP32_WALLET_FOLDER)
			if err := os.RemoveAll(pathDest); err != nil {
				fmt.Println("Clear dir failed", err)
			}
			if err := fileutils.MkdirAll(pathDest); err != nil {
				fmt.Println("Make dir failed", err)
			}
			components.PwdStr = WalletAccount.Password
			components.BipWallet, err = bip.ImportWalletFromMnemonicFile(filePath, pwd, WalletAccount.Password)
			if err != nil {
				fmt.Println("Import mnemonic failed", err)
			} else {
				b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
				components.BipAddress = base58.ChkEnc(b[:])
				OpenAndRunWallet()
			}
		}
	}
}

func selectAddress() {
	prompt := promptui.Select{
		Label: "Select wallet address",
		Items: components.OldAddresses,
	}

	_, result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input selection failed %v\n", err)
		_, result, err = prompt.Run()
	}

	components.XdagAddress = result
}
