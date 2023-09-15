package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-openapi/swag"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	"github.com/rivo/tview"
)

// fetch alerts data from alertmanager api
func (tui *TUI) getAlerts() {
	params := alert.NewGetAlertsParamsWithTimeout(10 * time.Second).WithContext(context.Background()).WithActive(swag.Bool(true)).WithSilenced(swag.Bool(false))
	tui.alerts(params)
}

// fetch filtered alerts data from alertmanager api
func (tui *TUI) getFilteredAlerts(filter []string) {
	params := alert.NewGetAlertsParamsWithTimeout(5 * time.Second).WithContext(context.Background()).WithFilter(filter).WithActive(swag.Bool(true)).WithSilenced(swag.Bool(false))
	tui.alerts(params)
}

func (tui *TUI) alerts(params *alert.GetAlertsParams) {
	err := tui.checkConn()
	if err != nil {
		tui.Errorf("%s", err)
		return
	}

	alerts, err := tui.amClient().Alert.GetAlerts(params)
	if err != nil {
		tui.Errorf("Error fetching alerts data: %s", err)
		return
	}

	tui.ClearPreviews()

	if len(alerts.Payload) == 0 {
		tui.Preview.SetText("[red]No matching alerts").SetTextAlign(tview.AlignCenter)
		return
	}

	tui.PreviewList.AddItem("Found "+strconv.Itoa(len(alerts.Payload))+" alerts 🔥", "", 0, nil)

	var mainText string
	var alertName string

	for _, a := range alerts.Payload {
		alertByte, err := json.MarshalIndent(a, "", "    ")
		if err != nil {
			log.Printf("Error marshaling alert: %s", err)
			continue
		}
		if a.Labels["severity"] != "" {
			switch a.Labels["severity"] {
			case "critical":
				alertName = "[red]" + a.Labels["alertname"]
			case "warning":
				alertName = "[yellow]" + a.Labels["alertname"]
			case "info":
				alertName = "[blue]" + a.Labels["alertname"]
			default:
				alertName = a.Labels["alertname"]
			}
		} else {
			alertName = a.Labels["alertname"]
		}
		if a.Annotations["description"] != "" {
			mainText = alertName + " - " + a.Annotations["description"]
		} else {
			mainText = alertName
		}
		tui.PreviewList.AddItem(mainText, fmt.Sprintf("[green]%s", string(alertByte)), 0, nil)
	}

	tui.PreviewList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		tui.Preview.Clear()
		tui.Preview.SetText(s2).SetTextAlign(tview.AlignLeft)
	})
}
