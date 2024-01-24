package components

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/fileutils"
	"goXdagWallet/i18n"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xlog"
	"image/color"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	cp "github.com/otiai10/copy"
)

type WalletState struct {
	HasAccount    bool
	Password      string
	WalletType    int // XDAG: 1, BIP32: 2
	MnemonicBytes []byte
}
type LogonWin struct {
	Win               fyne.Window
	BtnContainer      *fyne.Container
	LogonBtn          *widget.Button
	ProgressContainer *fyne.Container
	WalletState
}

var StatusInfo = canvas.NewText("", color.White)

func (l *LogonWin) NewLogonWindow(accountStatus int) {
	xlog.Info("Starting xdag wallet, version ", config.GetConfig().Version)
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
		getTestTitle())
	l.Win = w

	l.HasAccount = accountStatus >= 0

	//if accountStatus == HAS_BOTH || accountStatus == WALLET_NOT_FOUND {
	//	l.SelectWalletType()
	//} else {
	l.WalletType = accountStatus
	//}

	if l.HasAccount { // found wallet key file
		l.LogonBtn = widget.NewButton(i18n.GetString("LogonWindow_ConnectWallet"), l.connectClick)
	} else { // not fount
		l.LogonBtn = widget.NewButton(i18n.GetString("LogonWindow_RegisterWallet"), l.connectClick)
	}
	l.LogonBtn.Importance = widget.HighImportance

	//l.StatusInfo = canvas.NewText("", color.White)
	StatusInfo.Alignment = fyne.TextAlignCenter
	progress := widget.NewProgressBarInfinite()
	progress.Hide() // for primary color changing
	l.BtnContainer = container.New(layout.NewPaddedLayout(), l.LogonBtn)
	l.ProgressContainer = container.New(layout.NewPaddedLayout(), progress)

	settingBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		showLanguageDialog(i18n.GetString("SettingsWindow_ChooseLanguage"),
			i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), func(lang string) {
				config.GetConfig().CultureInfo = lang
				config.SaveConfig()
				i18n.LoadI18nStrings()
				l.Win.SetTitle(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
					getTestTitle())
				if accountStatus >= 0 { // found wallet key file
					l.LogonBtn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))
				} else if accountStatus == WALLET_NOT_FOUND { // not fount
					l.LogonBtn.SetText(i18n.GetString("LogonWindow_RegisterWallet"))
				}
			}, l.Win)
	})
	settingBtn.Resize(fyne.NewSize(20, 20))
	settingBtn.Importance = widget.HighImportance

	appearanceBtn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		fyneSetting := WalletApp.NewWindow(i18n.GetString("LogonWindow_AppearanceSetting"))
		fyneSetting.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		fyneSetting.Resize(fyne.NewSize(480, 480))
		fyneSetting.Show()
	})
	appearanceBtn.Resize(fyne.NewSize(20, 20))
	appearanceBtn.Importance = widget.HighImportance

	image := canvas.NewImageFromResource(resourceXj3Png)
	content := container.New(layout.NewMaxLayout(), image,
		container.New(layout.NewVBoxLayout(),
			container.New(layout.NewHBoxLayout(), layout.NewSpacer(), appearanceBtn, settingBtn),
			layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(),
			StatusInfo,
			l.BtnContainer,
			l.ProgressContainer,
			layout.NewSpacer()))
	w.SetContent(content)
	w.Resize(fyne.NewSize(435, 324))
	w.SetOnClosed(func() {
		xlog.CleanXdagLog(xlog.StdXdagLog)
		WalletApp.Quit()
		os.Exit(0)
	})
	w.CenterOnScreen()
	settingBtn.Refresh()
	appearanceBtn.Refresh()
}

func (l *LogonWin) StartConnect() {
	l.BtnContainer.Hide()
	l.ProgressContainer.Objects[0].Show() // for primary color changing
	StatusInfo.Text = i18n.GetString("LogonWindow_ConnectingAccount")
	canvas.Refresh(StatusInfo)
}
func (l *LogonWin) StartRegister() bool {
	l.BtnContainer.Hide()
	l.ProgressContainer.Objects[0].Show() // for primary color changing
	StatusInfo.Text = i18n.GetString("WalletState_Registering")
	canvas.Refresh(StatusInfo)

	w, b := NewBipWallet(PwdStr, 128)
	BipWallet = w
	return b

}

func (l *LogonWin) connectClick() {
	if config.GetConfig().Option.PoolAddress == "" {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("LogonWindow_NoPoolAddress"), l.Win)
		return
	}
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	accountStatus := l.WalletType
	if l.HasAccount {
		if accountStatus == HAS_BOTH {
			l.SelectWalletType()
		} else {
			l.showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
				i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
		}
	} else {
		l.CreateOrImport(pwd)
	}
}

func (l *LogonWin) showPasswordDialog(title, ok, dismiss string, parent fyne.Window) {
	l.LogonBtn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))
	wgt := widget.NewEntry()
	wgt.Password = true

	dialog.ShowCustomConfirm(title, ok, dismiss, wgt, func(b bool) {
		for i := range Password {
			Password[i] = 0
		}
		str := wgt.Text
		if b {
			if l.HasAccount {
				if len(str) > 0 {
					copy(Password[:], str)
					PwdStr = str
				}
				l.StartConnect()
				if l.WalletType == HAS_ONLY_XDAG {
					res := ConnectXdagWallet()
					if res == 0 {
						if XdagAddress == "" && len(OldAddresses) > 1 {
							showAddressSelect(l.Win)
						} else {
							NewWalletWindow(l.WalletType)
						}
					} else if res < -64 {
						l.loginIncorrect(i18n.GetString("LogonWindow_InitializeFailed"))
					} else {
						l.loginIncorrect(i18n.GetString("Message_PasswordIncorrect"))
					}
				} else if l.WalletType == HAS_ONLY_BIP {
					res := ConnectBipWallet(PwdStr)
					if res {
						NewWalletWindow(l.WalletType)
					} else {
						l.loginIncorrect(i18n.GetString("Message_PasswordIncorrect"))
					}
				}
			} else {
				l.Password = str
				l.ReShowPasswordDialog(i18n.GetString("PasswordWindow_RetypePassword"),
					i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
			}
		}

	}, parent)
}

func (l *LogonWin) loginIncorrect(msg string) {
	StatusInfo.Text = msg
	canvas.Refresh(StatusInfo)
	time.Sleep(time.Second * 5)
	l.Win.Close()
}

func (l *LogonWin) ReShowPasswordDialog(title, ok, dismiss string, parent fyne.Window) {
	wgt := widget.NewEntry()
	wgt.Password = true

	dialog.ShowCustomConfirm(title, ok, dismiss, wgt, func(b bool) {
		for i := range Password {
			Password[i] = 0
		}
		str := wgt.Text
		if b {
			if str == l.Password {
				if len(str) > 0 {
					copy(Password[:], str)
					PwdStr = str
				}
				if l.StartRegister() {
					NewWalletWindow(HAS_ONLY_BIP)
				} else {
					StatusInfo.Text = i18n.GetString("WalletState_Register_failed")
					canvas.Refresh(StatusInfo)
					time.Sleep(time.Second * 5)
					l.Win.Close()
				}

			} else {
				dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
					i18n.GetString("PasswordWindow_PasswordMismatch"), l.Win)
			}
		}
	}, parent)
}

func showLanguageDialog(title, ok, dismiss string, callback func(string), parent fyne.Window) {
	lang := "en-US"
	radio := widget.NewRadioGroup([]string{"English", "中文", "Français", "Русский", "Español", "Italiano"}, func(value string) {
		if value == "English" {
			lang = "en-US"
		} else if value == "Français" {
			lang = "fr-FR"
		} else if value == "Русский" {
			lang = "ru-RU"
		} else if value == "Italiano" {
			lang = "it-IT"
		} else if value == "Español" {
			lang = "sp-SP"
		} else {
			lang = "zh-CN"
		}
	})
	radio.Resize(fyne.Size{Width: radio.Size().Width, Height: radio.Size().Height * 2})
	dialog.ShowCustomConfirm(title, ok, dismiss, radio, func(b bool) {
		if b {
			callback(lang)
		}
	}, parent)
}

func GetAppIcon() fyne.Resource {
	return resourceWalletPng
}

func getTestTitle() string {
	var testNet string
	if config.GetConfig().Option.IsTestNet {
		if config.GetConfig().CultureInfo == "zh-CN" {
			testNet = "测试网"
		} else if config.GetConfig().CultureInfo == "fr-FR" {
			testNet = "Réseau Test"
		} else if config.GetConfig().CultureInfo == "ru-RU" {
			testNet = "тестовая сеть"
		} else {
			testNet = "Test Net"
		}
	}
	return testNet
}

func CheckOldWallet(walletDir string) bool {

	pathName := path.Join(walletDir, "dnet_key.dat")
	fi, err := os.Stat(pathName)
	if err != nil {
		return false
	}
	if fi.Size() != 2048 {
		return false
	}

	pathName = path.Join(walletDir, "wallet.dat")
	fi, err = os.Stat(pathName)
	if err != nil {
		return false
	}
	if fi.Size() != 32 {
		return false
	}

	pathName = path.Join(walletDir, "storage")
	_, err = os.Stat(pathName)

	return err == nil
}

func CopyOldWallet(walletDir string) error {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	pathDest := path.Join(pwd, "xdagj_dat")
	if err := os.RemoveAll(pathDest); err != nil {
		return err
	}
	if err := fileutils.MkdirAll(pathDest); err != nil {
		return err
	}
	if err := copyFile(walletDir, pathDest, "dnet_key.dat", 2048); err != nil {
		return err
	}
	if err := copyFile(walletDir, pathDest, "wallet.dat", 32); err != nil {
		return err
	}
	if err := cp.Copy(path.Join(walletDir, "storage"), path.Join(pathDest, "storage")); err != nil {
		return err
	}

	//l.importConfig(walletDir)

	return nil
}

func (l *LogonWin) ImportMnemonic(data []byte) error {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	pathDest := path.Join(pwd, common.BIP32_WALLET_FOLDER)
	if err := os.RemoveAll(pathDest); err != nil {
		return err
	}
	if err := fileutils.MkdirAll(pathDest); err != nil {
		return err
	}

	if len(data) < 12 || len(data) > 24 || len(data)%3 != 0 {
		xlog.Error("mnemonic file length error")
		return errors.New("mnemonic file length error")
	}

	destination, err := os.Create(path.Join(pathDest, common.BIP32_WALLET_FILE_NAME))
	if err != nil {
		xlog.Error(err)
		return err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, bytes.NewBuffer(data))
	if err != nil {
		xlog.Error(err)
		return err
	}
	if int(nBytes) != len(data) {
		xlog.Error("copy mnemonic file failed")
		return errors.New("copy mnemonic file failed")
	}
	return nil
}

func copyFile(walletDir, pathDest, fileName string, n int64) error {
	pathName := path.Join(walletDir, fileName)
	source, err := os.Open(pathName)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(path.Join(pathDest, fileName))
	if err != nil {
		return err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return err
	}
	if nBytes != n {
		return errors.New(fileName + " copy  failed")
	}
	return nil
}

func (l *LogonWin) importConfig(walletDir string) {

	configName := path.Join(walletDir, "wallet-config.json")
	_, err := os.Stat(configName)
	if err != nil {
		walletDir = strings.TrimSuffix(walletDir, `\`)
		walletDir = strings.TrimSuffix(walletDir, `/`)
		configName = path.Join(path.Dir(walletDir), "wallet-config.json")
		_, err = os.Stat(configName)
		if err != nil {
			return
		}
	}
	var oldConf config.Config
	data, err := os.ReadFile(configName)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &oldConf)
	if err != nil {
		return
	}

	if config.GetConfig().CultureInfo == "en-US" && oldConf.CultureInfo != "en-US" {

		config.GetConfig().CultureInfo = oldConf.CultureInfo

		i18n.LoadI18nStrings()

		l.Win.Resize(fyne.NewSize(410, 305))

		l.Win.SetTitle(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
			getTestTitle())

		l.LogonBtn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))

		l.Win.Resize(fyne.NewSize(410, 305))
	}

	config.GetConfig().Query = oldConf.Query

	for _, a := range oldConf.Addresses {
		config.InsertAddress(a)
	}
}

func (l *LogonWin) CreateOrImport(pwd string) {
	selectTypes := widget.NewRadioGroup([]string{
		i18n.GetString("LogonWindow_Create_Mnemonic"),
		i18n.GetString("LogonWindow_Import_Mnemonic"),
		i18n.GetString("LogonWindow_Import_NonMnemonic")},
		func(selected string) {
		})
	selectTypes.Selected = i18n.GetString("LogonWindow_Create_Mnemonic")
	selectTypes.Required = true

	content := []*widget.FormItem{
		{Text: i18n.GetString("LogonWindow_SelectWalletType") + ":", Widget: selectTypes},
	}
	query := dialog.NewForm(i18n.GetString("LogonWindow_SelectWalletTitle"),
		"   "+i18n.GetString("Common_Confirm")+"    ",
		"    "+i18n.GetString("Common_Cancel")+"     ",
		content,
		func(b bool) {
			if b {
				if selectTypes.Selected == i18n.GetString("LogonWindow_Import_NonMnemonic") {
					l.WalletType = HAS_ONLY_XDAG
					dlgOpen := dialog.NewFolderOpen(
						func(uri fyne.ListableURI, err error) {
							defer func() {
								l.Win.Resize(fyne.NewSize(410, 305))
							}()

							if uri == nil || err != nil {
								return
							}
							if !CheckOldWallet(uri.Path()) {
								dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
									i18n.GetString("WalletImport_WalletNotExist"), l.Win)
								return
							}
							if CopyOldWallet(uri.Path()) != nil {
								dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
									i18n.GetString("WalletImport_FilesCopyFailed"), l.Win)
								return
							}
							l.importConfig(uri.Path())
							l.HasAccount = true
							l.showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
								i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)

						}, l.Win)
					l.Win.Resize(fyne.NewSize(800, 500))
					dlgOpen.Resize(fyne.NewSize(800, 500))
					dlgOpen.Show()
				} else {
					l.WalletType = HAS_ONLY_BIP
					if selectTypes.Selected == i18n.GetString("LogonWindow_Create_Mnemonic") {
						pathDest := path.Join(pwd, common.BIP32_WALLET_FOLDER)
						if err := os.RemoveAll(pathDest); err != nil {
							dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
								i18n.GetString("WalletCreate_FilesFailed"), l.Win)
							return
						}
						if err := fileutils.MkdirAll(pathDest); err != nil {
							dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
								i18n.GetString("WalletCreate_FilesFailed"), l.Win)
							return
						}
						l.showPasswordDialog(i18n.GetString("PasswordWindow_SetPassword"),
							i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
					} else {
						dlgOpen := dialog.NewFileOpen(
							func(uri fyne.URIReadCloser, err error) {
								defer func() {
									l.Win.Resize(fyne.NewSize(410, 305))
								}()
								if uri == nil || err != nil {
									return
								}
								defer uri.Close()
								l.MnemonicBytes, err = io.ReadAll(uri)
								if err != nil || len(l.MnemonicBytes) == 0 {
									dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
										i18n.GetString("WalletImport_FilesCopyFailed"), l.Win)
									return
								}
								l.showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
									i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)

							}, l.Win)
						l.Win.Resize(fyne.NewSize(800, 500))
						dlgOpen.Resize(fyne.NewSize(800, 500))
						dlgOpen.Show()
					}
				}
			} else {
				l.Win.Close()
			}
		},
		l.Win)
	query.Resize(fyne.NewSize(150, 250))
	query.Show()
}

func (l *LogonWin) SelectWalletType() {
	selectTypes := widget.NewRadioGroup([]string{
		i18n.GetString("LogonWindow_WalletType_NonMnemonic"),
		i18n.GetString("LogonWindow_WalletType_Mnemonic")},
		func(selected string) {
		})
	selectTypes.Selected = i18n.GetString("LogonWindow_WalletType_NonMnemonic")
	selectTypes.Required = true

	content := []*widget.FormItem{
		{Text: i18n.GetString("LogonWindow_SelectWalletType") + ":", Widget: selectTypes},
	}
	query := dialog.NewForm(i18n.GetString("LogonWindow_SelectWalletTitle"),
		"   "+i18n.GetString("Common_Confirm")+"    ",
		"    "+i18n.GetString("Common_Cancel")+"     ",
		content,
		func(b bool) {
			if b {
				if selectTypes.Selected == i18n.GetString("LogonWindow_WalletType_NonMnemonic") {
					l.WalletType = HAS_ONLY_XDAG
				} else {
					l.WalletType = HAS_ONLY_BIP
				}
				l.showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
					i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
			} else {
				l.Win.Close()
			}
		},
		l.Win)
	query.Resize(fyne.NewSize(150, 200))
	query.Show()
}
