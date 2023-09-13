package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/prometheus/alertmanager/api/v2/client/silence"
	"github.com/rivo/tview"
)

// fetch silences data from alertmanager api
func (tui *TUI) getSilences() {
	err := tui.checkConn()
	if err != nil {
		tui.Errorf("%s", err)
		return
	}

	params := silence.NewGetSilencesParams().WithTimeout(5 * time.Second).WithContext(context.Background())
	silences, err := tui.amClient().Silence.GetSilences(params)
	if err != nil {
		tui.Errorf("Error fetching silences data: %s", err)
		return
	}

	tui.ClearPreviews()

	if len(silences.Payload) == 0 {
		tui.Preview.SetText("No silenced alerts 🔔").SetTextAlign(tview.AlignCenter)
		return
	}

	tui.PreviewList.SetTitle(" Silences ").SetTitleAlign(tview.AlignCenter)
	tui.PreviewList.AddItem("Total silences 🔕: "+strconv.Itoa(len(silences.Payload)), "", 0, nil)

	for _, silence := range silences.Payload {
		silenceByte, err := json.MarshalIndent(silence, "", "    ")
		if err != nil {
			log.Printf("Error marshaling silence: %s", err)
			continue
		}
		mainText := silence.EndsAt.String() + " - " + *silence.CreatedBy + " - " + *silence.Comment
		tui.PreviewList.AddItem(mainText, fmt.Sprintf("[green]%s", string(silenceByte)), 0, nil)
	}

	tui.PreviewList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		tui.Preview.Clear()
		tui.Preview.SetText(s2).SetTextAlign(tview.AlignLeft)
	})
}

func (tui *TUI) getFilteredSilences(filter []string) {
	err := tui.checkConn()
	if err != nil {
		tui.Errorf("%s", err)
		return
	}

	params := silence.NewGetSilencesParams().WithTimeout(5 * time.Second).WithContext(context.Background()).WithFilter(filter)
	silences, err := tui.amClient().Silence.GetSilences(params)
	if err != nil {
		tui.Errorf("Error fetching silences data: %s", err)
		return
	}

	tui.ClearPreviews()

	if len(silences.Payload) == 0 {
		tui.Preview.SetText("No silenced alerts 🔔").SetTextAlign(tview.AlignCenter)
		return
	}

	tui.PreviewList.SetTitle(" Silences ").SetTitleAlign(tview.AlignCenter)
	tui.PreviewList.AddItem("Total silences 🔕: "+strconv.Itoa(len(silences.Payload)), "", 0, nil)

	for _, silence := range silences.Payload {
		silenceByte, err := json.MarshalIndent(silence, "", "    ")
		if err != nil {
			log.Printf("Error marshaling silence: %s", err)
			continue
		}
		mainText := silence.EndsAt.String() + " - " + *silence.CreatedBy + " - " + *silence.Comment
		tui.PreviewList.AddItem(mainText, fmt.Sprintf("[green]%s", string(silenceByte)), 0, nil)
	}

	tui.PreviewList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		tui.Preview.Clear()
		tui.Preview.SetText(s2).SetTextAlign(tview.AlignLeft)
	})
}
