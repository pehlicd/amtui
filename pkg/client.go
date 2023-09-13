package pkg

import am "github.com/prometheus/alertmanager/api/v2/client"

const (
	BasePath = "/api/v2"
)

// create alertmanager client
func (tui *TUI) amClient() *am.AlertmanagerAPI {
	cfg := am.DefaultTransportConfig().WithHost(tui.Config.Host + ":" + tui.Config.Port).WithBasePath(BasePath).WithSchemes([]string{tui.Config.Scheme})
	return am.NewHTTPClientWithConfig(nil, cfg)
}
