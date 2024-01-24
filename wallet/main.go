// Package main provides various examples of Fyne API capabilities.
package main

import (
	"flag"
	"fmt"
	"goXdagWallet/cli"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/server"
	"goXdagWallet/xlog"
	"net"
	"os"
	"path"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var splashDone = make(chan struct{})

func init() {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)

	os.Setenv("FYNE_FONT", path.Join(pwd, "data", "myFont.ttf"))

}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func main() {
	xlog.SetLogFile("./", "go_wallet.log")
	config.InitConfig()
	accountStatus := components.Xdag_Wallet_fount() // search wallet data files

	mode := flag.String("mode", "gui", "run mode: gui, cli, server")
	ip := flag.String("ip", "127.0.0.1", "rpc server ip address")
	port := flag.Uint("port", 10001, "rpc server port number")
	flag.Parse()

	if *mode == "gui" {
		if err := i18n.LoadI18nStrings(); err != nil {
			xlog.Error(err)
			return
		}

		components.WalletApp = app.NewWithID("io.xdagj.wallet")
		components.WalletApp.SetIcon(components.GetAppIcon())
		if components.ShowSplashWindow(splashDone) {
			go func() {
				for range splashDone {
					components.LogonWindow.NewLogonWindow(accountStatus)
					components.LogonWindow.Win.Show()
					splashDone <- struct{}{}
				}
			}()
		} else {
			components.LogonWindow.NewLogonWindow(accountStatus)
			components.LogonWindow.Win.Show()
		}
		components.WalletApp.Run()
	} else if *mode == "cli" {
		cli.NewCli(accountStatus)
	} else if *mode == "server" {
		if accountStatus != components.HAS_ONLY_BIP {
			err := server.ImportServWallet()
			if err != nil {
				return
			}
		}
		if net.ParseIP(*ip) == nil {
			fmt.Println("Ip address format error")
			return
		}
		if *port < 1 || *port > 65535 {
			fmt.Println("Port number error")
			return
		}
		server.RunServer(*ip, strconv.Itoa(int(*port)))
	} else {
		fmt.Println("Unknown mode")
	}
	os.Unsetenv("FYNE_FONT")
}
