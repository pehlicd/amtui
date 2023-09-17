package pkg

import (
	clientruntime "github.com/go-openapi/runtime/client"
	am "github.com/prometheus/alertmanager/api/v2/client"
)

const (
	BasePath = "/api/v2"
)

// create alertmanager client
func (tui *TUI) amClient() *am.AlertmanagerAPI {
	address := tui.Config.Host + ":" + tui.Config.Port
	scheme := []string{tui.Config.Scheme}
	if tui.Config.Auth.Username != "" {
		cr := clientruntime.New(address, BasePath, scheme)
		cr.DefaultAuthentication = clientruntime.BasicAuth(tui.Config.Auth.Username, tui.Config.Auth.Password)
		c := am.New(cr, nil)
		return c
	} else {
		cfg := am.DefaultTransportConfig().WithHost(address).WithBasePath(BasePath).WithSchemes(scheme)
		return am.NewHTTPClientWithConfig(nil, cfg)
	}
}
