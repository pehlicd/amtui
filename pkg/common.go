package pkg

import (
	"fmt"
	"net"
	"time"

	"github.com/rivo/tview"
)

// checkConn dial tcp connection to alertmanager to be ensured if alertmanager server is up or not
func (tui *TUI) checkConn() error {
	conn, err := net.DialTimeout("tcp", tui.Config.Host+":"+tui.Config.Port, 1000*time.Millisecond)
	if err != nil {
		tui.Preview.Clear()
		return fmt.Errorf("error connecting to alertmanager host: %s", err)
	}
	defer conn.Close()
	return nil
}

// Errorf Create a function to print errors
func (tui *TUI) Errorf(format string, args ...interface{}) {
	tui.ClearPreviews()
	tui.Preview.SetText(fmt.Sprintf("[red]"+format, args...)).SetTextAlign(tview.AlignLeft)
}

// ClearPreviews Clear TUI previews
func (tui *TUI) ClearPreviews() {
	tui.PreviewList.Clear()
	tui.Preview.Clear()
}
