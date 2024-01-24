package components

import (
	"compress/gzip"
	"errors"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/buger/jsonparser"
)

type HistoryItem struct {
	Direction string
	Address   string
	Amount    string
	Time      string
	Remark    string
}

var historyStatus = binding.NewString()
var historyContainer *fyne.Container
var historyProgressContainer *fyne.Container
var historyRefreshContainer *fyne.Container
var historyNoResultContainer *fyne.Container
var historyData = make([]HistoryItem, 0, 10)
var historyTable *widget.Table
var curPage = 0
var pageCount = 0
var nextBtn *widget.Button
var prevBtn *widget.Button
var pageLabel = binding.NewString()
var queryParam string
var defaultParam string
var re = regexp.MustCompile(`^((19|20|21)[0-9][0-9])-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`)

func dateValidator() fyne.StringValidator {
	return func(text string) error {
		if text == "" {
			return nil
		}
		if re.MatchString(text) {
			return nil
		}
		return errors.New("not a valid date")
	}
}

func translateDirect(direction string) string {
	switch direction {
	case i18n.GetString("WalletWindow_History_Input"):
		return "input"
	case i18n.GetString("WalletWindow_History_Output"):
		return "output"
	default:
		return ""
	}
}

func makeQuery(amountFrom, amountTo, dateFrom, dateTo, remark, direction, timestamp string, remember bool) string {
	var condition string
	if len(dateFrom) > 0 {
		condition = condition + "&addresses_date_from=" + dateFrom
	}
	if len(dateTo) > 0 {
		condition = condition + "&addresses_date_to=" + dateTo
	}
	if len(amountFrom) > 0 {
		condition = condition + "&addresses_amount_from=" + amountFrom
	}
	if len(amountTo) > 0 {
		condition = condition + "&addresses_amount_to=" + amountTo
	}
	if len(remark) > 0 {
		condition = condition + "&addresses_remark=" + remark
	}
	if len(translateDirect(direction)) > 0 {
		condition = condition + "&addresses_directions[]=" + translateDirect(direction)
	}
	defaultParam = condition
	if remember {
		saveQuery(amountFrom, amountTo, remark, translateDirect(direction), timestamp)
	}
	return condition
}

func saveQuery(amountFrom, amountTo, remark, direction, timestamp string) {
	config.GetConfig().Query.AmountFrom = amountFrom
	config.GetConfig().Query.AmountTo = amountTo
	config.GetConfig().Query.Remark = remark
	config.GetConfig().Query.Direction = direction
	switch timestamp {
	case i18n.GetString("WalletWindow_Filter_Week"):
		config.GetConfig().Query.Timestamp = "week"
		break
	case i18n.GetString("WalletWindow_Filter_Month"):
		config.GetConfig().Query.Timestamp = "month"
		break
	case i18n.GetString("WalletWindow_Filter_Year"):
		config.GetConfig().Query.Timestamp = "year"
		break
	default:
		config.GetConfig().Query.Timestamp = "week"
	}
	config.SaveConfig()
}
func loadQuery() string {
	var condition string
	query := config.GetConfig().Query
	switch query.Timestamp {
	case "week":
		condition = condition + "&addresses_date_from=" +
			time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		break
	case "month":
		condition = condition + "&addresses_date_from=" +
			time.Now().AddDate(0, -1, 0).Format("2006-01-02")
		break
	case "year":
		condition = condition + "&addresses_date_from=" +
			time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		break
	default:
		condition = condition + "&addresses_date_from=" +
			time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}
	if len(query.Timestamp) > 0 {
		condition = condition + "&addresses_date_to=" + time.Now().Format("2006-01-02")
	}

	if len(query.AmountFrom) > 0 {
		condition = condition + "&addresses_amount_from=" + query.AmountFrom
	}
	if len(query.AmountTo) > 0 {
		condition = condition + "&addresses_amount_to=" + query.AmountTo
	}
	if len(query.Remark) > 0 {
		condition = condition + "&addresses_remark=" + query.Remark
	}
	if len(query.Direction) > 0 {
		condition = condition + "&addresses_directions[]=" + query.Direction
	}
	return condition
}
func setDefaultFilter(amountFrom, amountTo *numericalEntry, dateFrom, remark *widget.Entry,
	transferTime, direction *widget.RadioGroup) {
	query := config.GetConfig().Query
	switch query.Timestamp {
	case "week":
		transferTime.Selected = i18n.GetString("WalletWindow_Filter_Week")
		dateFrom.Text =
			time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		break
	case "month":
		transferTime.Selected = i18n.GetString("WalletWindow_Filter_Month")
		dateFrom.Text =
			time.Now().AddDate(0, -1, 0).Format("2006-01-02")
		break
	case "year":
		transferTime.Selected = i18n.GetString("WalletWindow_Filter_Year")
		dateFrom.Text =
			time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		break
	default:
		transferTime.Selected = i18n.GetString("WalletWindow_Filter_Week")
		dateFrom.Text =
			time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}

	if len(query.AmountFrom) > 0 {
		amountFrom.Text = query.AmountFrom
	}
	if len(query.AmountTo) > 0 {
		amountTo.Text = query.AmountTo
	}
	if len(query.Remark) > 0 {
		remark.Text = query.Remark
	}
	if query.Direction == "input" {
		direction.Selected = i18n.GetString("WalletWindow_History_Input")
	} else if query.Direction == "output" {
		direction.Selected = i18n.GetString("WalletWindow_History_Output")
	} else {
		direction.Selected = i18n.GetString("WalletWindow_Filter_AllDirect")
	}
}
func HistoryPage(w fyne.Window) *fyne.Container {
	defaultParam = loadQuery()
	queryParam = defaultParam
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		historyStatus.Set(i18n.GetString("WalletWindow_HistoryBusy"))
		historyProgressContainer.Show()
		historyRefreshContainer.Hide()
		historyData = historyData[:0]
		historyContainer.Remove(historyNoResultContainer)
		historyContainer.Remove(historyTable)
		go refreshTable(1, defaultParam)
	})
	refreshBtn.Importance = widget.HighImportance

	filterBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		dateFrom := widget.NewEntry()
		dateFrom.SetPlaceHolder(time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
		dateFrom.Validator = dateValidator()
		//dateFrom.Text = time.Now().AddDate(0, 0, -7).Format("2006-01-02")

		dateTo := widget.NewEntry()
		dateTo.SetPlaceHolder(time.Now().Format("2006-01-02"))
		dateTo.Validator = dateValidator()
		dateTo.Text = time.Now().Format("2006-01-02")

		remark := widget.NewEntry()
		amountFrom := newNumericalEntry()
		//amountFrom.Text = "0.1"
		amountTo := newNumericalEntry()
		direction := widget.NewRadioGroup([]string{
			i18n.GetString("WalletWindow_History_Input"),
			i18n.GetString("WalletWindow_History_Output"),
			i18n.GetString("WalletWindow_Filter_AllDirect")}, func(string) {
		})
		direction.Horizontal = true
		//direction.Selected = i18n.GetString("WalletWindow_Filter_AllDirect")
		direction.Required = true

		transferTime := widget.NewRadioGroup([]string{
			i18n.GetString("WalletWindow_Filter_Week"),
			i18n.GetString("WalletWindow_Filter_Month"),
			i18n.GetString("WalletWindow_Filter_Year")}, func(selected string) {
			dateTo.Text = time.Now().Format("2006-01-02")
			dateTo.Refresh()
			switch selected {
			case i18n.GetString("WalletWindow_Filter_Week"):
				dateFrom.Text = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
				dateFrom.Refresh()
				return
			case i18n.GetString("WalletWindow_Filter_Month"):
				dateFrom.Text = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
				dateFrom.Refresh()
				return
			case i18n.GetString("WalletWindow_Filter_Year"):
				dateFrom.Text = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
				dateFrom.Refresh()
				return
			}
		})
		transferTime.Horizontal = true
		//transferTime.Selected = i18n.GetString("WalletWindow_Filter_Week")
		transferTime.Required = true
		setDefaultFilter(amountFrom, amountTo, dateFrom, remark, transferTime, direction)
		remember := widget.NewCheck(i18n.GetString("WalletWindow_Filter_RemOption"),
			func(b bool) {})

		content := []*widget.FormItem{ // we can specify items in the constructor
			{Text: i18n.GetString("WalletWindow_Filter_AmountFrom") + ":", Widget: amountFrom},
			{Text: i18n.GetString("WalletWindow_Filter_AmountTo") + ":", Widget: amountTo},
			{Text: i18n.GetString("WalletWindow_HistoryColumns_TimeStamp") + ":", Widget: transferTime},
			{Text: i18n.GetString("WalletWindow_Filter_DateFrom") + ":", Widget: dateFrom},
			{Text: i18n.GetString("WalletWindow_Filter_DateTo") + ":", Widget: dateTo},
			{Text: i18n.GetString("WalletWindow_Transfer_Remark"), Widget: remark},
			{Text: i18n.GetString("WalletWindow_HistoryColumns_Direction") + ":", Widget: direction},
			{Text: i18n.GetString("WalletWindow_Filter_Remember"), Widget: remember},
		}

		query := dialog.NewForm(i18n.GetString("WalletWindow_History_Filter"),
			"   "+i18n.GetString("Common_Confirm")+"    ",
			"    "+i18n.GetString("Common_Cancel")+"     ",
			content,
			func(b bool) {
				if b {
					historyStatus.Set(i18n.GetString("WalletWindow_HistoryBusy"))
					historyProgressContainer.Show()
					historyRefreshContainer.Hide()
					historyData = historyData[:0]
					historyContainer.Remove(historyNoResultContainer)
					historyContainer.Remove(historyTable)
					queryParam = makeQuery(amountFrom.Text, amountTo.Text,
						dateFrom.Text, dateTo.Text, remark.Text, direction.Selected,
						transferTime.Selected, remember.Checked)
					go refreshTable(1, queryParam)
				}
			},
			w)
		query.Resize(fyne.NewSize(150, 200))
		query.Show()
	})
	filterBtn.Importance = widget.HighImportance

	nextBtn = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		if curPage < pageCount {
			curPage += 1
			historyProgressContainer.Show()
			historyRefreshContainer.Hide()
			historyData = historyData[:0]
			historyContainer.Remove(historyTable)
			go refreshTable(curPage, queryParam)
		}
	})
	nextBtn.Importance = widget.HighImportance

	prevBtn = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		if curPage > 1 {
			curPage -= 1
			historyProgressContainer.Show()
			historyRefreshContainer.Hide()
			historyData = historyData[:0]
			historyContainer.Remove(historyNoResultContainer)
			historyContainer.Remove(historyTable)
			go refreshTable(curPage, queryParam)
		}
	})
	prevBtn.Importance = widget.HighImportance
	label := widget.NewLabelWithData(pageLabel)

	historyProgressContainer = container.New(layout.NewPaddedLayout(), widget.NewProgressBarInfinite())
	historyRefreshContainer = container.New(layout.NewHBoxLayout(),
		prevBtn, label, nextBtn, layout.NewSpacer(), filterBtn, refreshBtn)
	historyRefreshContainer.Hide()
	historyContainer = container.New(layout.NewMaxLayout())
	HistoryTitle := widget.NewLabelWithData(historyStatus)
	historyStatus.Set(i18n.GetString("WalletWindow_HistoryBusy"))
	pageLabel.Set("1/1")

	historyNoResultContainer = container.NewHBox(layout.NewSpacer(),
		widget.NewLabel(i18n.GetString("WalletWindow_Filter_NoResult")),
		layout.NewSpacer())
	go refreshTable(1, defaultParam)
	top := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), HistoryTitle, layout.NewSpacer()),
		historyRefreshContainer,
		historyProgressContainer)
	return container.New(
		layout.NewBorderLayout(top, nil, nil, nil),
		top,
		historyContainer)
}

func refreshTable(page int, query string) {
	var body []byte
	if getApiUrl() == "" {
		return
	}
	var address string
	if LogonWindow.WalletType == HAS_ONLY_XDAG {
		address = XdagAddress
	} else {
		address = BipAddress
	}
	err := getUrl(getApiUrl(), address, query, page, &body)
	if err != nil {
		historyProgressContainer.Hide()
		historyRefreshContainer.Show()
		historyStatus.Set(i18n.GetString("WalletWindow_HistoryError"))
		return
	}

	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		direction, _ := jsonparser.GetString(value, "direction")
		address, _ := jsonparser.GetString(value, "address")
		amount, _ := jsonparser.GetString(value, "amount")
		times, _ := jsonparser.GetString(value, "time")
		remark, _ := jsonparser.GetString(value, "remark")
		historyData = append(historyData, HistoryItem{
			Direction: direction,
			Address:   address,
			Amount:    amount,
			Time:      times,
			Remark:    remark,
		})
	}, "block_as_address")
	historyProgressContainer.Hide()
	historyRefreshContainer.Show()
	if err != nil {
		historyStatus.Set(i18n.GetString("WalletWindow_HistoryError"))
		return
	}
	total, _ := jsonparser.GetInt(body, "addresses_pagination", "total")
	lastPage, _ := jsonparser.GetInt(body, "addresses_pagination", "last_page")
	pageCount = int(lastPage)

	current, _ := jsonparser.GetInt(body, "addresses_pagination", "current_page")
	prev, _ := jsonparser.GetString(body, "addresses_pagination", "links", "prev")
	if len(prev) == 0 {
		prevBtn.Disable()
	} else {
		prevBtn.Enable()
	}
	curPage = int(current)

	next, _ := jsonparser.GetString(body, "addresses_pagination", "links", "next")
	if len(next) == 0 {
		nextBtn.Disable()
	} else {
		nextBtn.Enable()
	}
	pageLabel.Set(strconv.Itoa(curPage) + "/" + strconv.Itoa(pageCount))

	historyStatus.Set(i18n.GetString("WalletWindow_HistoryColumns_BlockAddress") + " : " + address)

	historyTable = widget.NewTable(
		func() (int, int) { return len(historyData) + 1, 5 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_Direction"))
				case 1:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_Amount"))
				case 2:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_PartnerAddress"))
				case 3:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_TimeStamp"))
				case 4:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_Remark"))
				default:
					label.SetText("cell")
				}
			} else {
				switch id.Col {
				case 0:
					if historyData[id.Row-1].Direction == "input" {
						label.SetText(i18n.GetString("WalletWindow_History_Input"))
					} else {
						label.SetText(i18n.GetString("WalletWindow_History_Output"))
					}

				case 1:
					label.SetText(historyData[id.Row-1].Amount)
				case 2:
					label.SetText(historyData[id.Row-1].Address)
				case 3:
					label.SetText(historyData[id.Row-1].Time)
				case 4:
					label.SetText(historyData[id.Row-1].Remark)
				default:
					label.SetText("cell")
				}
			}

		})
	if config.GetConfig().CultureInfo == "ru-RU" ||
		config.GetConfig().CultureInfo == "fr-FR" {
		historyTable.SetColumnWidth(0, 122)
	} else {
		historyTable.SetColumnWidth(0, 82)
	}

	historyTable.SetColumnWidth(1, 178)
	historyTable.SetColumnWidth(2, 372)
	historyTable.SetColumnWidth(3, 188)
	historyTable.SetColumnWidth(4, 152)
	historyTable.Refresh()
	if total == 0 {
		historyContainer.Add(historyNoResultContainer)
	} else {
		historyContainer.Add(historyTable)
	}

}
func getUrl(apiUrl, address, query string, page int, body *[]byte) error {
	urlString := apiUrl + "/" + address +
		"?addresses_per_page=10&addresses_page=" + strconv.Itoa(page)
	if len(query) > 0 {
		urlString = urlString + query
	}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return err
	}

	// 表单方式(必须)
	//req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Apache-HttpClient/4.3.1")

	client := &http.Client{
		Timeout: 8 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err == nil {
			defer reader.Close()
		} else {
			return err
		}
	default:
		reader = resp.Body
	}
	*body, err = io.ReadAll(reader)
	if err != nil {
		return err
	}
	return nil
}

func getApiUrl() string {
	if config.GetConfig().Option.IsTestNet {
		if apiUrl, err := url.Parse(config.GetConfig().Option.TestnetApiUrl); err == nil && apiUrl.IsAbs() {
			return config.GetConfig().Option.TestnetApiUrl
		} else {
			return "https://testexplorer.xdag.io/api/block"
		}
	} else {
		if len(config.GetConfig().Option.TestnetApiUrl) > 0 {
			return config.GetConfig().Option.TestnetApiUrl
		} else {
			return "https://explorer.xdag.io/api/block"
		}
	}
}
