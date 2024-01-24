package components

import (
	"fmt"
	"goXdagWallet/i18n"
	"goXdagWallet/xlog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var DonaTransStatus = widget.NewLabel("")
var DonaTransBtnContainer *fyne.Container
var DonaTransProgressContainer *fyne.Container
var DonaAddressEntry = widget.NewEntry()
var coinImages = map[int]*canvas.Image{
	1:   canvas.NewImageFromResource(resourceCoin01Png),
	5:   canvas.NewImageFromResource(resourceCoin05Png),
	10:  canvas.NewImageFromResource(resourceCoin10Png),
	20:  canvas.NewImageFromResource(resourceCoin20Png),
	50:  canvas.NewImageFromResource(resourceCoin50Png),
	100: canvas.NewImageFromResource(resourceCoin100Png),
}

func DonatePage(w fyne.Window) *fyne.Container {
	amount := newNumericalEntry()
	remark := widget.NewEntry()
	DonaAddressEntry.Validator = addressValidator()
	remark.Validator = remarkValidator()

	DonaTransProgressContainer = container.New(layout.NewPaddedLayout(), widget.NewProgressBarInfinite())
	DonaTransProgressContainer.Hide()

	btn := widget.NewButtonWithIcon(i18n.GetString("WalletWindow_TabDonate"), theme.ConfirmIcon(),
		func() {
			showPwdConfirm(w, func() {
				fromAccountPrivKey, fromAddress, fromValue := SelTransFromAddr()
				if !checkInput(fromValue, DonaAddressEntry.Text, amount.Text, remark.Text, fromAddress, w) {
					return
				}

				message := fmt.Sprintf(i18n.GetString("TransferWindow_ConfirmTransfer"), amount.Text, DonaAddressEntry.Text)
				//fmt.Println(message)
				dialog.ShowConfirm(i18n.GetString("Common_ConfirmTitle"),
					message, func(b bool) {
						if b {
							DonaTransProgressContainer.Show()
							DonaTransBtnContainer.Hide()
							TransProgressContainer.Show()
							TransBtnContainer.Hide()
							DonaTransStatus.Text = i18n.GetString("TransferWindow_CommittingTransaction")
							TransStatus.Text = i18n.GetString("TransferWindow_CommittingTransaction")
							_, err := TransferRpc(fromAddress, DonaAddressEntry.Text, amount.Text, remark.Text, fromAccountPrivKey)
							if err == nil {
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
	DonaTransBtnContainer = container.New(layout.NewPaddedLayout(), btn)

	setDonateValue := func(value string) {
		amount.SetText(value)
	}
	top := container.NewVBox(
		container.NewHBox(layout.NewSpacer(),
			newCoinBtn(1, func() { setDonateValue("1") }), layout.NewSpacer(),
			newCoinBtn(5, func() { setDonateValue("5") }), layout.NewSpacer(),
			newCoinBtn(10, func() { setDonateValue("10") }), layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			newCoinBtn(20, func() { setDonateValue("20") }), layout.NewSpacer(),
			newCoinBtn(50, func() { setDonateValue("50") }), layout.NewSpacer(),
			newCoinBtn(100, func() { setDonateValue("100") }), layout.NewSpacer()),
		container.NewHBox(
			layout.NewSpacer(), DonaTransStatus, layout.NewSpacer()),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: i18n.GetString("WalletWindow_Transfer_ToAddress"), Widget: DonaAddressEntry},
				{Text: i18n.GetString("WalletWindow_Transfer_Amount"), Widget: amount},
				{Text: i18n.GetString("WalletWindow_Transfer_Remark"), Widget: remark},
			},
		}),
		container.NewVBox(
			DonaTransProgressContainer,
			DonaTransBtnContainer),
		widget.NewLabel(""),
	)
	DonaAddressEntry.Text = CommunityAddress
	DonaAddressEntry.Disable()

	return container.New(
		layout.NewBorderLayout(top, nil, nil, nil),
		top,
	)
}

func newCoinBtn(value int, tap func()) *fyne.Container {
	image := coinImages[value] // this is an image.Image
	// image := canvas.NewImageFromFile("./images/coin" + value + ".png") // this is an image.Image
	image.SetMinSize(fyne.NewSize(100, 100))
	// image.FillMode = canvas.ImageFillContain
	// This button will be placed "over" the image using a Padded box
	openButton := widget.NewButton("", tap)
	openButton.Resize(fyne.NewSize(100, 100))
	// this encapsulate the image and the button
	box := container.NewMax(image, openButton)
	box.Resize(fyne.NewSize(100, 100))
	return box
}
