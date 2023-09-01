# AMTUI - Alertmanager Terminal User Interface

```

 █████╗ ███╗   ███╗████████╗██╗   ██╗██╗
██╔══██╗████╗ ████║╚══██╔══╝██║   ██║██║
███████║██╔████╔██║   ██║   ██║   ██║██║
██╔══██║██║╚██╔╝██║   ██║   ██║   ██║██║
██║  ██║██║ ╚═╝ ██║   ██║   ╚██████╔╝██║
╚═╝  ╚═╝╚═╝     ╚═╝   ╚═╝    ╚═════╝ ╚═╝
                             
```

AMTUI is a terminal-based user interface (TUI) application that allows you to interact with Prometheus Alertmanager using your terminal. It provides a convenient way to monitor alerts, view silences, and check the status of Alertmanager instances right from your command line.

## Features

- View active alerts with details such as severity, alert name, and description.
- Browse and review existing silences in Alertmanager.
- Check the general status of your Alertmanager instance.

## Installation

To use AMTUI, you'll need to have Go installed on your system. Then, you can install AMTUI using the following steps:

1. Clone the repository:

```bash
git clone https://github.com/pehlicd/amtui.git
```

2. Navigate to the project directory:

```bash
cd amtui
```

3. Build the application:

```bash
go build
```

4. Run the application:

```bash
./amtui
```

## Usage

Once you've launched AMTUI, you can navigate through different sections using the following keyboard shortcuts:

- Press `1` to view and interact with active alerts.
- Press `2` to see existing silences.
- Press `3` to check the general status of your Alertmanager instance.

### Keyboard Shortcuts

- `q`: Quit the application.
- `l`: Focus on the preview list.
- `h`: Focus on the sidebar list.
- `j`: Move focus to the preview.
- `k`: Move focus to the preview list.
- `ESC`: Return focus to the sidebar list.

## Configuration

AMTUI uses a configuration file to connect to your Alertmanager instance. By default, the application will look for a configuration file at `~/.amtui.yaml`. If the configuration file doesn't exist, AMTUI will guide you through creating it with the necessary connection details.

You can also specify connection details using command-line flags:

```bash
./amtui -host your-alertmanager-host -port your-alertmanager-port -scheme http
```

## Dependencies

AMTUI uses the following dependencies:

- `github.com/gdamore/tcell/v2`: Terminal handling and screen painting.
- `github.com/prometheus/alertmanager/api/v2/client`: Alertmanager API client.
- `github.com/rivo/tview`: Terminal-based interactive viewer.
- `gopkg.in/yaml.v3`: YAML support for the configuration file.

## Contributing

If you'd like to contribute to AMTUI, feel free to submit pull requests or open issues on the [GitHub repository](https://github.com/pehlicd/amtui). Your feedback and contributions are highly appreciated!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Developed by [Furkan Pehlivan](https://github.com/pehlicd) - [Project Repository](https://github.com/pehlicd/amtui)