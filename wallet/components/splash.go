package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"time"
)

func ShowSplashWindow(done chan struct{}) bool {
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		w := drv.CreateSplashWindow()
		img := canvas.NewImageFromResource(resourceSplashPng)
		w.SetContent(container.New(layout.NewMaxLayout(), img))
		w.Resize(fyne.NewSize(570, 380))
		w.CenterOnScreen()
		w.Show()
		go func() {
			<-time.After(time.Second * 2)
			done <- struct{}{}
			<-done
			w.Close()
			close(done)
		}()
		return true
	}
	return false
}
