package components

import (
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xdago/cryptography"
	xdagoUtils "goXdagWallet/xdago/utils"
	bip "goXdagWallet/xdago/wallet"
	"goXdagWallet/xlog"
	"os"
	"path"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/utf8string"
)

var chanBalance = make(chan int, 1)

func Xdag_Wallet_fount() int {
	hasXdagWallet := 0
	hasBip32Wallet := 0
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	pathName := path.Join(pwd, "xdagj_dat", "dnet_key.dat")
	// change current working directory
	os.Chdir(pwd)

	fi, err := os.Stat(pathName)
	if err != nil {
		hasXdagWallet = -1
	} else if fi.Size() != 2048 {
		hasXdagWallet = -1
	}
	pathName = path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME)
	//os.Chdir(pwd)

	fi, err = os.Stat(pathName)
	if err != nil {
		hasBip32Wallet = -1
	} else if fi.Size() < 126 { // file is 125 bytes when wallet data without Mnemonic
		hasBip32Wallet = -1
	}
	if hasXdagWallet == -1 && hasBip32Wallet == -1 {
		return WALLET_NOT_FOUND // no wallet
	}
	if hasXdagWallet == 0 && hasBip32Wallet == 0 { // has both wallets
		return HAS_BOTH
	} else if hasXdagWallet == 0 { // only has xdag wallet
		return HAS_ONLY_XDAG
	} else {
		return HAS_ONLY_BIP // only has bip32_bip44 wallet
	}
}

func ConnectBipWallet(password string) bool {
	xlog.Info("Initializing cryptography...")
	xlog.Info("Reading wallet...")
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	wallet := bip.NewWallet(path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	res := wallet.UnlockWallet(password)
	if wallet.IsHdWalletInitialized() {
		xlog.Info("Reading Mnemonic...")
	}
	if res && wallet.IsHdWalletInitialized() {
		BipWallet = &wallet
		return true
	}
	return false

}
func NewBipWallet(password string, bitSize int) (*bip.Wallet, bool) {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	if len(LogonWindow.MnemonicBytes) > 0 {
		xlog.Info("import Mnemonic...")
		fmt.Println("Importing Mnemonic...")
		wallet, err := bip.ImportWalletFromMnemonicStr(string(LogonWindow.MnemonicBytes), pwd, PwdStr)
		if err != nil {
			xlog.Error(err)
			return nil, false
		} else {
			return wallet, true
		}
	} else {
		xlog.Info("creating Mnemonic...")
		fmt.Printf("Creating Mnemonic...")
		wallet := bip.NewWallet(path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
		wallet.UnlockWallet(password)
		wallet.InitializeHdWallet(bip.NewMnemonic(bitSize))
		wallet.AddAccountWithNextHdKey()
		res := wallet.Flush()
		return &wallet, res
	}
}

func ValidateXdagAddress(address string) bool {
	_, err := xdagoUtils.Address2Hash(address)
	return err == nil
}

func ValidateBipAddress(address string) bool {
	_, err := checkBase58Address(address)
	return err == nil
}

func ValidateRemark(remark string) bool {
	return utf8string.NewString(remark).IsASCII() && len(remark) < 32
}

func NewWalletWindow(walletExists int) {
	if WalletWindow != nil {
		return
	}

	if walletExists != HAS_ONLY_XDAG {
		b := cryptography.ToBytesAddress(BipWallet.GetDefKey())
		BipAddress = base58.ChkEnc(b[:])
	}

	address, balance := getBalance()
	LogonWindow.Win.Hide()
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
		getTestTitle())
	WalletWindow = w
	w.SetMaster()
	LogonWindow.Win.Content().Resize(fyne.NewSize(0, 0))
	tabs := container.NewAppTabs()
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAccount"),
		theme.HomeIcon(), AccountPage(address, balance, WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabTransfer"),
		theme.MailSendIcon(), TransferPage(WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabHistory"),
		theme.ContentPasteIcon(), HistoryPage(WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabDonate"),
		theme.NewThemedResource(donateIconRes), DonatePage(WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAbout"),
		theme.InfoIcon(), AboutPage(WalletWindow)))

	if fyne.CurrentDevice().IsMobile() {
		tabs.SetTabLocation(container.TabLocationBottom)
	} else {
		tabs.SetTabLocation(container.TabLocationLeading)
	}

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(640, 480))
	w.CenterOnScreen()
	go checkBalance()
	w.SetOnClosed(func() {
		xlog.CleanXdagLog(xlog.StdXdagLog)
		chanBalance <- 1
		WalletApp.Quit()
		os.Exit(0)
	})
	w.Show()
}

func checkBalance() {
	for {
		select {
		case <-chanBalance:
			return
		case <-time.After(time.Second * 130):
			getBalance()
		}
	}
}

func getBalance() (string, string) {
	if LogonWindow.WalletType == HAS_ONLY_BIP {
		balance, err := BalanceRpc(BipAddress)
		if err == nil {
			if balance != "" {
				BipBalance = balance
				AccountBalance.Set(balance)
				return BipAddress, BipBalance
			} else {
				xlog.Error("get bip32 account balance error.")
				return BipAddress, ""
			}
		} else {
			xlog.Error(err)
			return BipAddress, ""
		}
	} else { // LogonWindow.WalletType == HAS_ONLY_XDAG
		balance, err := BalanceRpc(XdagAddress)
		if err == nil {
			if balance != "" {
				XdagBalance = balance
				AccountBalance.Set(balance)
				return XdagAddress, XdagBalance
			} else {
				xlog.Error("get xdag account balance error.")
				return XdagAddress, ""
			}
		} else {
			xlog.Error(err)
			return XdagAddress, ""
		}
	}
}

func showAddressSelect(w fyne.Window) {
	radio := widget.NewRadioGroup(OldAddresses, func(value string) {
		XdagAddress = value
	})
	radio.SetSelected(OldAddresses[0])
	XdagAddress = OldAddresses[0]
	dialog.ShowCustomConfirm(i18n.GetString("PasswordWindow_RetypePassword"),
		i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), radio, func(b bool) {
			NewWalletWindow(HAS_ONLY_XDAG)
		}, w)
}
