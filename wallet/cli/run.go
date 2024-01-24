package cli

import (
	"errors"
	"fmt"
	"goXdagWallet/components"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/cryptography"
	"goXdagWallet/xdago/secp256k1"
	xdagoUtils "goXdagWallet/xdago/utils"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
)

var spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

func OpenAndRunWallet() {
	fmt.Println("Initializing cryptography...")
	fmt.Println("Reading wallet...")
	spin.Start()
	if WalletAccount.WalletType == components.HAS_ONLY_XDAG {
		res := components.ConnectXdagWallet()
		spin.Stop()
		if res == 0 {
			if components.XdagAddress == "" && len(components.OldAddresses) > 0 {
				selectAddress()
			}
			RunWallet(WalletAccount.WalletType)
		} else {
			fmt.Println("Password incorrect")
		}
	} else if WalletAccount.WalletType == components.HAS_ONLY_BIP {
		res := components.ConnectBipWallet(WalletAccount.Password)
		spin.Stop()
		if res {
			b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
			components.BipAddress = base58.ChkEnc(b[:])
			RunWallet(WalletAccount.WalletType)
		} else {
			fmt.Println("Password incorrect")
		}
	}
}

var validateCmd = func(input string) error {
	if input != "help" && input != "exit" && input != "account" &&
		input != "balance" && !strings.HasPrefix(input, "xfer ") &&
		input != "mnemonic" && !strings.HasPrefix(input, "export ") {
		return errors.New("unknown command, input 'help' to list available commands")
	}
	if strings.HasPrefix(input, "xfer ") {
		items := strings.Fields(input)
		if len(items) != 4 {
			return errors.New("transfer command parameters error")
		}
		if !components.ValidateBipAddress(items[2]) {
			return errors.New("address format error")
		}
		value, err := strconv.ParseFloat(items[1], 64)
		if err != nil || value <= 0.0 {
			return errors.New("amount number error")
		}
		if !components.ValidateRemark(items[3]) {
			return errors.New("remark format error")
		}
	}

	if strings.HasPrefix(input, "export ") {
		items := strings.Fields(input)
		if len(items) != 2 {
			return errors.New("export command parameters error")
		}
		if strings.HasSuffix(items[1], `\`) || strings.HasSuffix(items[1], `/`) {
			return errors.New("path to export file error")
		}
		if strings.Contains(items[1], `\`) || strings.Contains(items[1], `/`) {
			folder, _ := path.Split(items[1])
			if !xdagoUtils.FileExists(folder) {
				return errors.New("path to export file not exist")
			}
		}
	}
	return nil
}

func RunWallet(walletExists int) {
	prompt := promptui.Prompt{
		Label:    "Command > ",
		Validate: validateCmd,
	}
	for {
		result, err := prompt.Run()

		for err != nil {
			fmt.Printf("Input command failed %v\n", err)
			result, err = prompt.Run()
		}
		switch result {
		case "help":
			fmt.Println("---------------------------------------------------------")
			fmt.Println("      help -- display commands list")
			fmt.Println("      exit -- exit cli wallet")
			fmt.Println("   account -- display address of wallet account")
			fmt.Println("   balance -- display balance of wallet account")
			fmt.Println("xfer V A R -- transfer V coins to address A with remark R")
			if walletExists == components.HAS_ONLY_BIP {
				fmt.Println("  mnemonic -- display mnemonic of wallet account")
				fmt.Println("  export P -- export mnemonic to file P")
			}
			fmt.Println("---------------------------------------------------------")
		case "exit":
			os.Exit(0)
		case "account":
			if walletExists == components.HAS_ONLY_BIP {
				fmt.Println(components.BipAddress)
			} else if walletExists == components.HAS_ONLY_XDAG {
				fmt.Println(components.XdagAddress)
			}
		case "balance":
			var balance string
			var errBlc error
			spin.Start()
			if walletExists == components.HAS_ONLY_BIP {
				balance, errBlc = components.BalanceRpc(components.BipAddress)
			} else if walletExists == components.HAS_ONLY_XDAG {
				balance, errBlc = components.BalanceRpc(components.XdagAddress)
			}
			spin.Stop()
			if errBlc != nil {
				fmt.Println("Get balance failed", errBlc)
			} else {
				fmt.Println(balance)
			}
		case "mnemonic":
			if walletExists == components.HAS_ONLY_BIP {
				fmt.Println(components.BipWallet.GetMnemonic())
			} else {
				fmt.Println("It's a Non Mnemonic wallet")
			}
		}

		if strings.HasPrefix(result, "xfer ") {
			var fromAddress string
			var fromKey *secp256k1.PrivateKey
			items := strings.Fields(result)
			spin.Start()
			if walletExists == components.HAS_ONLY_BIP {
				fromAddress = components.BipAddress
				fromKey = components.BipWallet.GetDefKey()
			} else if walletExists == components.HAS_ONLY_XDAG {
				fromAddress = components.XdagAddress
				fromKey = components.XdagKey
			}

			fromValue, errBlc := components.BalanceRpc(fromAddress)
			if errBlc != nil {
				spin.Stop()
				fmt.Println("Get balance failed", errBlc)
				continue
			}

			if !checkInput(fromValue, items[2], items[1], items[3], fromAddress) {
				spin.Stop()
				continue
			}
			_, errTx := components.TransferRpc(fromAddress, items[2], items[1], items[3], fromKey)
			spin.Stop()
			if errTx != nil {
				fmt.Println("Transfer failed", err)
			} else {
				fmt.Println("Transfer has been committed. Please wait for a while to get it done")
			}
		}

		if strings.HasPrefix(result, "export ") {
			if walletExists == components.HAS_ONLY_BIP {
				items := strings.Fields(result)
				errExp := components.BipWallet.ExportMnemonic(items[1])
				if errExp != nil {
					fmt.Println("Export mnemonic failed", err)
				} else {
					fmt.Println("Export mnemonic success")
				}
			} else {
				fmt.Println("It's a Non Mnemonic wallet")
			}
		}
	}
}

func checkInput(fromValue, toAddr, amount, remark, fromAddress string) bool {
	if len(toAddr) == 0 || !components.ValidateBipAddress(toAddr) || fromAddress == toAddr {
		fmt.Println("Receive Address format is incorrect.")
		return false
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil || value <= 0.0 {
		fmt.Println("Amount should be a positive number.")
		return false
	}

	balance, _ := strconv.ParseFloat(fromValue, 64)
	if balance < value {
		fmt.Println("Insufficient amount")
		return false
	}

	if len(remark) > 0 && !components.ValidateRemark(remark) {
		fmt.Println("Remark format is incorrect")
		return false
	}
	return true
}
