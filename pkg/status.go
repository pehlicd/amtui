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
	err := tui.checkConn()
	if err != nil {
		tui.Errorf("%s", err)
		return
	}

	params := general.NewGetStatusParams().WithTimeout(10 * time.Second).WithContext(context.Background())
	status, err := tui.amClient().General.GetStatus(params)
	if err != nil {
		tui.Errorf("Error fetching status data: %s", err)
	}

	tui.ClearPreviews()

	statusByte, err := json.MarshalIndent(status.Payload, "", "    ")
	if err != nil {
		tui.Errorf("Error marshaling status: %s", err)
	}

	tui.PreviewList.SetTitle(" Status ").SetTitleAlign(tview.AlignCenter)
	tui.Preview.SetText(fmt.Sprintf("[green]%s", string(statusByte))).SetTextAlign(tview.AlignLeft)
}
