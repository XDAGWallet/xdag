package components

import (
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	CommunityOldAddress = "FQglVQtb60vQv2DOWEUL7yh3smtj7g1s"
	CommunityAddress    = "4duPWMbYUgAifVYkKDCWxLvRRkSByf5gb"
)

func AboutPage(w fyne.Window) *fyne.Container {
	link, _ := url.Parse("https://xdag.io/")
	tele, _ := url.Parse("https://t.me/dagger_cryptocurrency")
	discord, _ := url.Parse("https://discord.gg/YxXUVQJ")

	var testNet string
	if config.GetConfig().Option.IsTestNet {
		if config.GetConfig().CultureInfo == "zh-CN" {
			testNet = "测试网"
		} else if config.GetConfig().CultureInfo == "fr-FR" {
			testNet = "Réseau d'essai"
		} else if config.GetConfig().CultureInfo == "ru-RU" {
			testNet = "тестовая сеть"
		} else if config.GetConfig().CultureInfo == "it-IT" {
			testNet = "Rete di prova"
		} else if config.GetConfig().CultureInfo == "sp-SP" {
			testNet = "Red de prueba"
		} else {
			testNet = "Test Net"
		}
	}

	cnContainer := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("XDAG钱包（"+config.GetConfig().Version+"）"+testNet), layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("XDAG是基于PoW共识算法和DAG技术的加密货币，解决了传统区块链技术"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("的瓶颈问题，是首个可挖矿的DAG项目。"),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("社区网站："), widget.NewHyperlink("xdag.io", link),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("加入社区："), widget.NewHyperlink("Discord", discord),
			widget.NewHyperlink("Telegram", tele),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("为社区团队捐赠XDAG:"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel(CommunityAddress),
			widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
				w.Clipboard().SetContent(CommunityAddress)
				dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
					i18n.GetString("WalletWindow_AddressCopied"), w)
			}),
			layout.NewSpacer()),
	)
	enContainer := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("XDAG wallet("+config.GetConfig().Version+") "+testNet), layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("XDAG is a novel application of Directed Acyclic Graph (DAG) technology that"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("solves the issues currently facing blockchain technology."),
			layout.NewSpacer()), container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("The first mineable DAG."),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Website："), widget.NewHyperlink("xdag.io", link),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Join XDAG："), widget.NewHyperlink("Discord", discord),
			widget.NewHyperlink("Telegram", tele),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Donate XDAG to Community Team:"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel(CommunityAddress),
			widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
				w.Clipboard().SetContent(CommunityAddress)
				dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
					i18n.GetString("WalletWindow_AddressCopied"), w)
			}),
			layout.NewSpacer()),
	)
	frContainer := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Portefeuille XDAG ("+config.GetConfig().Version+") "+testNet), layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("XDAG est une nouvelle application de la technologie des Graphes Acycliques Dirigés (DAG) qui"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("résout les problèmes auxquels est actuellement confrontée la technologie blockchain."),
			layout.NewSpacer()), container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Le premier DAG minable."),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Site web："), widget.NewHyperlink("xdag.io", link),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Communauté XDAG："), widget.NewHyperlink("Discord", discord),
			widget.NewHyperlink("Telegram", tele),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Faire un don à l équipe communautaire XDAG:"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel(CommunityAddress),
			widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
				w.Clipboard().SetContent(CommunityAddress)
				dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
					i18n.GetString("WalletWindow_AddressCopied"), w)
			}),
			layout.NewSpacer()),
	)

	ruContainer := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("XDAG кошелек ("+config.GetConfig().Version+") "+testNet), layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("XDAG - новое применение технологии направленного ациклического графа (DAG),"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("решающей проблемы, с которыми в настоящее время сталкивается технология блокчейн."),
			layout.NewSpacer()), container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Первый DAG, который можно майнить."),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Веб-сайт："), widget.NewHyperlink("xdag.io", link),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Присоединиться к XDAG："), widget.NewHyperlink("Discord", discord),
			widget.NewHyperlink("Telegram", tele),
			layout.NewSpacer()),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel("Подарить XDAG команде сообщества:"),
			layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(),
			widget.NewLabel(CommunityAddress),
			widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
				w.Clipboard().SetContent(CommunityAddress)
				dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
					i18n.GetString("WalletWindow_AddressCopied"), w)
			}),
			layout.NewSpacer()),
	)

	if config.GetConfig().CultureInfo == "zh-CN" {
		return cnContainer
	} else if config.GetConfig().CultureInfo == "fr-FR" {
		return frContainer
	} else if config.GetConfig().CultureInfo == "ru-RU" {
		return ruContainer
	} else {
		return enContainer
	}

}
