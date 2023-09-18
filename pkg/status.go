package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prometheus/alertmanager/api/v2/client/general"
	"github.com/rivo/tview"
)

// fetch status data from alertmanager api
func (tui *TUI) getStatus() {
	params := general.NewGetStatusParams().WithTimeout(5 * time.Second).WithContext(context.Background())
	status, err := tui.amClient().General.GetStatus(params)
	if err != nil {
		tui.Errorf("Error fetching status data: %s", err)
		return
	}

	tui.ClearPreviews()

	statusByte, err := json.MarshalIndent(status.Payload, "", "    ")
	if err != nil {
		tui.Errorf("Error marshaling status: %s", err)
	}

	tui.PreviewList.SetTitle(" Status ").SetTitleAlign(tview.AlignCenter)

	tui.PreviewList.AddItem("Status", "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" Uptime: %s", status.Payload.Uptime), "", 0, nil)
	tui.PreviewList.AddItem("*", "", 0, nil)
	tui.PreviewList.AddItem("Cluster Status", "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" Status: %s", *status.Payload.Cluster.Status), "", 0, nil)
	tui.PreviewList.AddItem(" Peers:", "", 0, nil)
	if len(status.Payload.Cluster.Peers) > 0 {
		for _, peer := range status.Payload.Cluster.Peers {
			tui.PreviewList.AddItem(fmt.Sprintf("   - Name: %s Address: %s", *peer.Name, *peer.Address), "", 0, nil)
		}
	}
	tui.PreviewList.AddItem("*", "", 0, nil)
	tui.PreviewList.AddItem("Version Information", "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" Branch: %s", *status.Payload.VersionInfo.Branch), "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" Build Date: %s", *status.Payload.VersionInfo.BuildDate), "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" Build User: %s", *status.Payload.VersionInfo.BuildUser), "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" GoVersion: %s", *status.Payload.VersionInfo.GoVersion), "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" Revision: %s", *status.Payload.VersionInfo.Revision), "", 0, nil)
	tui.PreviewList.AddItem(fmt.Sprintf(" Version: %s", *status.Payload.VersionInfo.Version), "", 0, nil)

	tui.Preview.SetText(fmt.Sprintf("[green]%s", string(statusByte))).SetTextAlign(tview.AlignLeft)
}
