package pkg

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	BasePath        = "/api/v2"
	TitleFooterView = "AMTUI - Alertmanager TUI Client\ngithub.com/pehlicd/amtui"
)

type TUI struct {
	App         *tview.Application
	SidebarList *tview.List
	PreviewList *tview.List
	Preview     *tview.TextView
	Grid        *tview.Grid
	FooterText  *tview.TextView
	Filter      *tview.InputField
	Config      Config
}

func InitTUI() *TUI {
	tui := TUI{App: tview.NewApplication()}

	tui.SidebarList = tview.NewList().ShowSecondaryText(false)
	tui.PreviewList = tview.NewList().ShowSecondaryText(false).SetSelectedBackgroundColor(tcell.ColorIndigo).SetSelectedTextColor(tcell.ColorWhite)
	tui.Preview = tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetScrollable(true)
	tui.Filter = tview.NewInputField().SetLabel("Filter: ").SetFieldBackgroundColor(tcell.ColorIndigo).SetLabelColor(tcell.ColorWhite).SetFieldTextColor(tcell.ColorWhite).SetDoneFunc(func(key tcell.Key) {
		// check if Alerts option is selected from SidebarList or not
		if tui.SidebarList.GetCurrentItem() != 0 {
			tui.ClearPreviews()
			tui.Preview.SetText("[red]Please select Alerts option from Navigation")
			return
		}
		// if search field is empty, return all alerts
		if tui.Filter.GetText() == "" {
			tui.getAlerts()
			return
		}
		// if search field is not empty, return alerts based on search field
		tui.PreviewList.Clear()
		filter := strings.Split(tui.Filter.GetText(), ",")
		tui.getFilteredAlerts(filter)
		tui.App.SetFocus(tui.PreviewList)
	}).SetPlaceholder("Custom matcher, e.g. env=\"production\"").SetPlaceholderTextColor(tcell.ColorIndigo)
	tui.FooterText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(TitleFooterView).SetTextColor(tcell.ColorGray).SetWordWrap(true)

	tui.PreviewList.SetTitle("").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	tui.SidebarList.SetTitle(" Navigation ").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	tui.SidebarList.AddItem("Alerts", "", '1', tui.getAlerts)
	tui.SidebarList.AddItem("Silences", "", '2', tui.getSilences)
	tui.SidebarList.AddItem("Status", "", '3', tui.getStatus)
	tui.Preview.SetTitle("").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	tui.Filter.SetTitle(" Filter ").SetTitleAlign(tview.AlignCenter).SetBorder(true)

	tui.Grid = tview.NewGrid().
		SetRows(3, 0, 0, 2).
		SetColumns(20, 0).
		AddItem(tui.SidebarList, 0, 0, 3, 1, 0, 0, true).
		AddItem(tui.Filter, 0, 1, 1, 1, 0, 0, false).
		AddItem(tui.PreviewList, 1, 1, 1, 1, 0, 0, false).
		AddItem(tui.Preview, 2, 1, 1, 1, 0, 0, false).
		AddItem(tui.FooterText, 3, 0, 1, 2, 0, 0, false)

	// configuration management
	tui.Config = initConfig()

	// listen for keyboard events and if q pressed, exit if l pressed in SidebarList focus on PreviewList if h is pressed in PreviewList focus on SidebarList
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				tui.App.Stop()
				return nil
			case 'l':
				if tui.App.GetFocus() == tui.SidebarList {
					tui.App.SetFocus(tui.PreviewList)
				}
				return nil
			case 'h':
				tui.App.SetFocus(tui.SidebarList)
				return nil
			case 'j':
				if tui.App.GetFocus() == tui.PreviewList {
					tui.App.SetFocus(tui.Preview)
				}
				return nil
			case 'k':
				if tui.App.GetFocus() == tui.Preview {
					tui.App.SetFocus(tui.PreviewList)
				}
				return nil
			}
		} else if event.Key() == tcell.KeyEsc {
			if tui.App.GetFocus() == tui.Filter {
				tui.App.SetFocus(tui.PreviewList)
				return nil
			}
			tui.App.SetFocus(tui.SidebarList)
			return nil
		} else if event.Key() == tcell.KeyCtrlF {
			tui.App.SetFocus(tui.Filter)
			return nil
		}
		return event
	})
	return &tui
}

func (tui *TUI) Start() error {
	return tui.App.SetRoot(tui.Grid, true).Run()
}
