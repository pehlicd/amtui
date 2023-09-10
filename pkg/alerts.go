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
	err := tui.checkConn()
	if err != nil {
		tui.Errorf("%s", err)
		return
	}

	params := alert.NewGetAlertsParamsWithTimeout(10 * time.Second).WithContext(context.Background()).WithActive(swag.Bool(true)).WithSilenced(swag.Bool(false))
	alerts, err := tui.amClient().Alert.GetAlerts(params)
	if err != nil {
		tui.Errorf("Error fetching alerts data: %s", err)
		return
	}

	tui.ClearPreviews()
	tui.PreviewList.SetTitle(" Alerts ").SetTitleAlign(tview.AlignCenter).SetBorder(true)

	if len(alerts.Payload) == 0 {
		tui.Preview.SetText("[green]No alerts 🎉").SetTextAlign(tview.AlignCenter)
		return
	}

	tui.PreviewList.AddItem("Total active alerts 🔥: "+strconv.Itoa(len(alerts.Payload)), "", 0, nil)

	var mainText string
	var alertName string

	for _, alert := range alerts.Payload {
		alertByte, err := json.MarshalIndent(alert, "", "    ")
		if err != nil {
			log.Printf("Error marshaling alert: %s", err)
			continue
		}
		if alert.Labels["severity"] != "" {
			switch alert.Labels["severity"] {
			case "critical":
				alertName = "[red]" + alert.Labels["alertname"]
			case "warning":
				alertName = "[yellow]" + alert.Labels["alertname"]
			case "info":
				alertName = "[blue]" + alert.Labels["alertname"]
			default:
				alertName = alert.Labels["alertname"]
			}
		} else {
			alertName = alert.Labels["alertname"]
		}
		if alert.Annotations["description"] != "" {
			mainText = alertName + " - " + alert.Annotations["description"]
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

// fetch filtered alerts data from alertmanager api
func (tui *TUI) getFilteredAlerts(filter []string) {
	err := tui.checkConn()
	if err != nil {
		tui.Errorf("%s", err)
		return
	}

	params := alert.NewGetAlertsParamsWithTimeout(10 * time.Second).WithContext(context.Background()).WithFilter(filter).WithActive(swag.Bool(true)).WithSilenced(swag.Bool(false))
	alerts, err := tui.amClient().Alert.GetAlerts(params)
	if err != nil {
		tui.Errorf("Error fetching alerts data: %s", err)
		return
	}

	if len(alerts.Payload) == 0 {
		tui.Preview.SetText("[red]No matching alerts").SetTextAlign(tview.AlignCenter)
		return
	}

	tui.PreviewList.AddItem("Found "+strconv.Itoa(len(alerts.Payload))+" alerts 🔥", "", 0, nil)

	var mainText string
	var alertName string

	for _, alert := range alerts.Payload {
		alertByte, err := json.MarshalIndent(alert, "", "    ")
		if err != nil {
			log.Printf("Error marshaling alert: %s", err)
			continue
		}
		if alert.Labels["severity"] != "" {
			switch alert.Labels["severity"] {
			case "critical":
				alertName = "[red]" + alert.Labels["alertname"]
			case "warning":
				alertName = "[yellow]" + alert.Labels["alertname"]
			case "info":
				alertName = "[blue]" + alert.Labels["alertname"]
			default:
				alertName = alert.Labels["alertname"]
			}
		} else {
			alertName = alert.Labels["alertname"]
		}
		if alert.Annotations["description"] != "" {
			mainText = alertName + " - " + alert.Annotations["description"]
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
