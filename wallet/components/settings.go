package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
)

func SettingsPage(w fyne.Window) fyne.CanvasObject {
	s := settings.NewSettings()
	return s.LoadAppearanceScreen(w)
}
