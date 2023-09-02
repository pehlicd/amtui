/*

 █████╗ ███╗   ███╗████████╗██╗   ██╗██╗
██╔══██╗████╗ ████║╚══██╔══╝██║   ██║██║
███████║██╔████╔██║   ██║   ██║   ██║██║
██╔══██║██║╚██╔╝██║   ██║   ██║   ██║██║
██║  ██║██║ ╚═╝ ██║   ██║   ╚██████╔╝██║
╚═╝  ╚═╝╚═╝     ╚═╝   ╚═╝    ╚═════╝ ╚═╝

*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/go-openapi/swag"
	am "github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	"github.com/prometheus/alertmanager/api/v2/client/silence"
	"github.com/rivo/tview"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// API URLs
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
	Config      Config
}

type Config struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Scheme string `yaml:"scheme"`
}

func initConfig() Config {
	// Initialize Viper
	viper.SetConfigName(".amtui")          // Configuration file name without extension
	viper.SetConfigType("yaml")            // Configuration file type
	viper.AddConfigPath(os.Getenv("HOME")) // Search for the configuration file in the $HOME directory

	// Set default values for your configuration struct
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "9093")
	viper.SetDefault("scheme", "http")

	var config Config

	// Allow command-line flags to override the configuration
	flag.StringVar(&config.Host, "host", config.Host, "Alertmanager host")
	flag.StringVar(&config.Port, "port", config.Port, "Alertmanager port")
	flag.StringVar(&config.Scheme, "scheme", config.Scheme, "Alertmanager scheme http or https is supported")
	flag.Parse()

	// Bind environment variables (optional)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("AMTUI")

	//if flags are set, overwrite config file
	if config.Host != "" && config.Port != "" && config.Scheme != "" {
		viper.Set("host", config.Host)
		viper.Set("port", config.Port)
		viper.Set("scheme", config.Scheme)
		err := viper.WriteConfig()
		if err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// Handle errors when the configuration file is not found or is invalid
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults.")
			// Write the default configuration to a new file
			if err := viper.SafeWriteConfig(); err != nil {
				log.Fatalf("Error creating config file: %v", err)
			}
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	// Merge flags into the configuration
	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		log.Fatalf("Error binding flags: %v", err)
	}

	// Unmarshal the configuration into your Config struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return config
}

func tuiInit() *TUI {
	tui := TUI{App: tview.NewApplication()}
	tui.SidebarList = tview.NewList().ShowSecondaryText(false)
	tui.PreviewList = tview.NewList().ShowSecondaryText(false).SetSelectedBackgroundColor(tcell.ColorBrown)
	tui.Preview = tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetScrollable(true)
	tui.FooterText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(TitleFooterView).SetTextColor(tcell.ColorGray)
	tui.PreviewList.SetTitle("").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	tui.SidebarList.SetTitle(" Navigation ").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	tui.SidebarList.AddItem("Alerts", "", '1', tui.getAlerts)
	tui.SidebarList.AddItem("Silences", "", '2', tui.silences)
	tui.SidebarList.AddItem("Status", "", '3', tui.status)
	tui.Preview.SetTitle("").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	tui.Grid = tview.NewGrid().
		SetRows(0, 0, 3).
		SetColumns(20, 0).
		AddItem(tui.SidebarList, 0, 0, 2, 1, 0, 0, true).
		AddItem(tui.PreviewList, 0, 1, 1, 1, 0, 0, false).
		AddItem(tui.Preview, 1, 1, 1, 1, 0, 0, false).
		AddItem(tui.FooterText, 2, 0, 1, 2, 0, 0, false)
	// check if config file exists
	tui.Config = initConfig()
	return &tui
}

func main() {
	tui := tuiInit()

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
			tui.App.SetFocus(tui.SidebarList)
			return nil
		}
		return event
	})

	if err := tui.Start(); err != nil {
		fmt.Printf("Error running app: %s", err)
		os.Exit(1)
	}
}

func (tui *TUI) Start() error {
	return tui.App.SetRoot(tui.Grid, true).Run()
}

// create alertmanager client
func (tui *TUI) amClient() *am.AlertmanagerAPI {
	cfg := am.DefaultTransportConfig().WithHost(tui.Config.Host + ":" + tui.Config.Port).WithBasePath(BasePath).WithSchemes([]string{tui.Config.Scheme})
	return am.NewHTTPClientWithConfig(nil, cfg)
}

// fetch alerts data from alertmanager api
func (tui *TUI) getAlerts() {
	err := tui.checkConn()
	if err != nil {
		tui.Preview.SetText(fmt.Sprintf("%s", err))
		return
	}

	alerts, err := tui.amClient().Alert.GetAlerts(&alert.GetAlertsParams{Silenced: swag.Bool(false), Active: swag.Bool(true), Context: context.Background()})
	if err != nil {
		tui.Preview.SetText(fmt.Sprintf("Error fetching alerts data: %s", err))
		return
	}
	tui.PreviewList.Clear()
	tui.Preview.Clear()
	if len(alerts.Payload) == 0 {
		tui.Preview.SetText("[green]No alerts 🎉").SetTextAlign(tview.AlignCenter)
		return
	}
	tui.PreviewList.SetTitle(" Alerts ").SetTitleAlign(tview.AlignCenter).SetBorder(true)
	var mainText string
	var alertName string
	tui.PreviewList.AddItem("[red]Total active alerts 🔥: [white]"+strconv.Itoa(len(alerts.Payload)), "", 0, nil)
	for _, alert := range alerts.Payload {
		alertByte, err := json.MarshalIndent(alert, "", "    ")
		if err != nil {
			fmt.Printf("Error marshaling alert: %s", err)
		}
		if alert.Labels["severity"] != "" {
			switch alert.Labels["severity"] {
			case "critical":
				alertName = "[red]" + alert.Labels["alertname"]
			case "warning":
				alertName = "[yellow]" + alert.Labels["alertname"]
			case "info":
				alertName = "[blue]" + alert.Labels["alertname"]
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

// fetch silences data from alertmanager api
func (tui *TUI) silences() {
	err := tui.checkConn()
	if err != nil {
		tui.Preview.SetText(fmt.Sprintf("%s", err))
		return
	}
	silences, err := tui.amClient().Silence.GetSilences(&silence.GetSilencesParams{Context: context.Background()})
	if err != nil {
		tui.Preview.SetText(fmt.Sprintf("Error fetching silences data: %s", err))
		return
	}
	tui.Preview.Clear()
	tui.PreviewList.Clear()
	if len(silences.Payload) == 0 {
		tui.Preview.SetText("No silenced alerts 🔔").SetTextAlign(tview.AlignCenter)
		return
	}
	tui.PreviewList.SetTitle(" Silences ").SetTitleAlign(tview.AlignCenter)
	tui.PreviewList.AddItem("Total silences 🔕: "+strconv.Itoa(len(silences.Payload)), "", 0, nil)
	for _, silence := range silences.Payload {
		silenceByte, err := json.MarshalIndent(silence, "", "    ")
		if err != nil {
			fmt.Printf("Error marshaling silence: %s", err)
			continue
		}
		mainText := silence.EndsAt.String() + " - " + *silence.CreatedBy + " - " + *silence.Comment
		tui.PreviewList.AddItem(mainText, fmt.Sprintf("[white]%s", string(silenceByte)), 0, nil)
	}

	tui.PreviewList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		tui.Preview.Clear()
		tui.Preview.SetText(s2).SetTextAlign(tview.AlignLeft)
	})
}

// fetch status data from alertmanager api
func (tui *TUI) status() {
	err := tui.checkConn()
	if err != nil {
		tui.Preview.SetText(fmt.Sprintf("%s", err))
		return
	}
	status, err := tui.amClient().General.GetStatus(nil)
	if err != nil {
		tui.Preview.SetText(fmt.Sprintf("Error fetching status data: %s", err))
	}
	tui.Preview.Clear()
	tui.PreviewList.Clear()
	statusByte, err := json.MarshalIndent(status.Payload, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling status: %s", err)
	}
	tui.PreviewList.SetTitle("Status").SetTitleAlign(tview.AlignCenter)

	tui.Preview.SetText(fmt.Sprintf("[white]%s", string(statusByte))).SetTextAlign(tview.AlignLeft)
}

// send http get request to alertmanager api to be ensure if it is up or not
func (tui *TUI) checkConn() error {
	resp, err := http.Get(tui.Config.Scheme + "://" + tui.Config.Host + ":" + tui.Config.Port + BasePath + "/status")
	if err != nil {
		tui.Preview.Clear()
		return fmt.Errorf("error connecting to alertmanager api: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		tui.Preview.Clear()
		return fmt.Errorf("error connecting to alertmanager api: %s", resp.Status)
	}
	return nil
}
