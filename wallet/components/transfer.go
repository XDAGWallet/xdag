package components

import (
	"errors"
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xdago/secp256k1"
	"goXdagWallet/xlog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var TransStatus = widget.NewLabel("")
var TransBtnContainer *fyne.Container
var TransProgressContainer *fyne.Container
var AddressList *widget.List
var AddressEntry = widget.NewEntry()

//var SelectedAddress *widget.RadioGroup

func addressValidator() fyne.StringValidator {
	return func(text string) error {
		if text == "" {
			return nil
		}
		if ValidateBipAddress(text) {
			return nil
		}
		return errors.New(i18n.GetString("TransferWindow_AccountFormatError"))
	}
}

func remarkValidator() fyne.StringValidator {
	return func(text string) error {
		if text == "" {
			return nil
		}
		if ValidateRemark(text) {
			return nil
		}
		return errors.New(i18n.GetString("TransferWindow_RemarkFormatError"))
	}
}

func TransferPage(w fyne.Window) *fyne.Container {
	amount := newNumericalEntry()
	remark := widget.NewEntry()
	AddressEntry.Validator = addressValidator()
	remark.Validator = remarkValidator()

	TransProgressContainer = container.New(layout.NewPaddedLayout(), widget.NewProgressBarInfinite())
	TransProgressContainer.Hide()

	//if LogonWindow.WalletType == HAS_BOTH {
	//	SelectedAddress = widget.NewRadioGroup([]string{
	//		XdagAddress,
	//		BipAddress}, func(selected string) {
	//	})
	//	SelectedAddress.Selected = XdagAddress
	//	//SelectedAddress.Horizontal = true
	//	SelectedAddress.Required = true
	//
	//}
	btn := widget.NewButtonWithIcon(i18n.GetString("TransferWindow_TransferTitle"), theme.ConfirmIcon(),
		func() {
			showPwdConfirm(w, func() {

				fromAccountPrivKey, fromAddress, fromValue := SelTransFromAddr()

				if !checkInput(fromValue, AddressEntry.Text, amount.Text, remark.Text, fromAddress, w) {
					return
				}
				if fromAddress == AddressEntry.Text {
					return
				}

				message := fmt.Sprintf(i18n.GetString("TransferWindow_ConfirmTransfer"), amount.Text, AddressEntry.Text)
				//fmt.Println(message)
				dialog.ShowConfirm(i18n.GetString("Common_ConfirmTitle"),
					message, func(b bool) {
						if b {
							TransProgressContainer.Show()
							TransBtnContainer.Hide()
							DonaTransProgressContainer.Show()
							DonaTransBtnContainer.Hide()
							TransStatus.Text = i18n.GetString("TransferWindow_CommittingTransaction")
							DonaTransStatus.Text = i18n.GetString("TransferWindow_CommittingTransaction")
							_, err := TransferRpc(fromAddress, AddressEntry.Text, amount.Text, remark.Text, fromAccountPrivKey)
							if err == nil {
								config.InsertAddress(AddressEntry.Text)
								setTransferDone()
							} else {
								xlog.Error(err)
								setTransferError(err.Error())
							}
						}
					}, w)
			})
		})
	btn.Importance = widget.HighImportance
	TransBtnContainer = container.New(layout.NewPaddedLayout(), btn)
	makeAddrList()

	top := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(), TransStatus, layout.NewSpacer()),
		//makeTopBar(),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: i18n.GetString("WalletWindow_Transfer_ToAddress"), Widget: AddressEntry},
				{Text: i18n.GetString("WalletWindow_Transfer_Amount"), Widget: amount},
				{Text: i18n.GetString("WalletWindow_Transfer_Remark"), Widget: remark},
			},
		}),
		container.NewVBox(
			TransProgressContainer,
			TransBtnContainer),
		widget.NewLabel(""),
		widget.NewLabel(i18n.GetString("TransferWindow_MostRecently")),
	)

	return container.New(
		layout.NewBorderLayout(top, nil, nil, nil),
		top,
		AddressList,
	)
}

func checkInput(fromValue, toAddr, amount, remark, fromAddress string, window fyne.Window) bool {
	if len(toAddr) == 0 || !ValidateBipAddress(toAddr) || fromAddress == toAddr {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_AccountFormatError"), window)
		return false
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil || value <= 0.0 {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_AmountFormatError"), window)
		return false
	}

	balance, _ := strconv.ParseFloat(fromValue, 64)
	if balance < value {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_InsufficientAmount"), window)
		return false
	}

	if len(remark) > 0 && !ValidateRemark(remark) {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_RemarkFormatError"), window)
		return false
	}
	return true
}

func makeAddrList() {
	AddressList = widget.NewList(
		func() int {
			return len(config.GetConfig().Addresses)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.ContentCopyIcon()), widget.NewLabel("address"),
				layout.NewSpacer(), widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {

				}))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			address := config.GetConfig().Addresses[id]
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(address)
			item.(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
				dialog.ShowConfirm(i18n.GetString("Common_ConfirmTitle"),
					fmt.Sprintf(i18n.GetString("TransferWindow_ConfirmDelete"), address),
					func(b bool) {
						if b {
							config.DeleteAddress(id)
							AddressList.Refresh()
						}
					}, WalletWindow)

			}
		})
	AddressList.OnSelected = func(id widget.ListItemID) {
		AddressEntry.SetText(config.GetConfig().Addresses[id])
	}
}

func setTransferDone() {
	TransProgressContainer.Hide()
	TransBtnContainer.Show()
	DonaTransProgressContainer.Hide()
	DonaTransBtnContainer.Show()
	dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
		i18n.GetString("TransferWindow_CommitSuccess"), WalletWindow)
	makeAddrList()
}

func setTransferError(e string) {
	TransProgressContainer.Hide()
	TransBtnContainer.Show()
	DonaTransProgressContainer.Hide()
	DonaTransBtnContainer.Show()
	dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
		i18n.GetString("TransferWindow_CommitFailed")+e, WalletWindow)
}

func SelTransFromAddr() (*secp256k1.PrivateKey, string, string) {
	if LogonWindow.WalletType == HAS_ONLY_XDAG {
		return XdagKey, XdagAddress, XdagBalance
	} else { //  LogonWindow.WalletType == HAS_ONLY_BIP
		return BipWallet.GetDefKey(), BipAddress, BipBalance
	}
	//else { // WalletType == HAS_BOTH
	//	if SelectedAddress.Selected == XdagAddress {
	//		return XdagKey, XdagAddress, XdagBalance
	//	} else {
	//		return BipWallet.GetDefKey(), BipAddress, BipBalance
	//	}
	//}
}

//
//func makeTopBar() fyne.CanvasObject {
//	if LogonWindow.WalletType == HAS_BOTH {
//		return container.NewVBox(
//			container.NewHBox(layout.NewSpacer(), TransStatus, layout.NewSpacer()),
//			container.NewHBox(widget.NewLabel(i18n.GetString("TransferWindow_FromAddress")),
//				SelectedAddress))
//	} else {
//		return container.NewHBox(
//			layout.NewSpacer(), TransStatus, layout.NewSpacer())
//	}
//}
