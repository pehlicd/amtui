# AMTUI - Alertmanager Terminal User Interface

![Go version](https://img.shields.io/github/go-mod/go-version/pehlicd/amtui) ![Release](https://img.shields.io/github/v/release/pehlicd/amtui) [![Go Report Card](https://goreportcard.com/badge/github.com/pehlicd/amtui)](https://goreportcard.com/report/github.com/pehlicd/amtui) ![License](https://img.shields.io/github/license/pehlicd/amtui) ![Discord](https://img.shields.io/discord/1152242693557202975?logo=discord&link=https%3A%2F%2Fdiscord.gg%2FJXnFz5j42n)

```

 █████╗ ███╗   ███╗████████╗██╗   ██╗██╗
██╔══██╗████╗ ████║╚══██╔══╝██║   ██║██║
███████║██╔████╔██║   ██║   ██║   ██║██║
██╔══██║██║╚██╔╝██║   ██║   ██║   ██║██║
██║  ██║██║ ╚═╝ ██║   ██║   ╚██████╔╝██║
╚═╝  ╚═╝╚═╝     ╚═╝   ╚═╝    ╚═════╝ ╚═╝
                             
```

AMTUI is a terminal-based user interface (TUI) application that allows you to interact with Prometheus Alertmanager using your terminal. It provides a convenient way to monitor alerts, view silences, and check the status of Alertmanager instances right from your command line.

<p align="center">
    <img src="./static/demo.gif" alt="AMTUI Demo"/>
</p>

## Features

- View active alerts with details such as severity, alert name, and description.
- Browse and review existing silences in Alertmanager.
- Filter alerts and silences using matchers.
- Check the general status of your Alertmanager instance.

## Installation

### Using Homebrew
You can install AMTUI using the [Homebrew](https://brew.sh/) package manager:

```bash
brew tap pehlicd/tap
brew install amtui
```

### Using go install
You can install AMTUI using the `go install` command:

```bash
go install github.com/pehlicd/amtui@latest
```

### From Releases
You can download the latest release of AMTUI from the [GitHub releases page](https://github.com/pehlicd/amtui/releases).

### From Source
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
- `CTRL + F`: Focus on the filter input.
- `ESC`: Return focus to the sidebar list.

## Configuration

AMTUI uses a configuration file to connect to your Alertmanager instance. By default, the application will look for a configuration file at `~/.amtui.yaml`. If the configuration file doesn't exist, AMTUI will guide you through creating it with the necessary connection details.

You can also specify connection details using command-line flags:

```bash
amtui --host 127.0.0.1 --port 9093 --scheme http
```

AMTUI also supports basic authentication. You can specify the username and password using the `--username` and `--password` flags:

```bash
amtui --host 127.0.0.1 --port 9093 --scheme http --username admin --password admin
```

## Dependencies

AMTUI uses the following dependencies:

- `github.com/gdamore/tcell/v2`: Terminal handling and screen painting.
- `github.com/prometheus/alertmanager/api/v2/client`: Alertmanager API client.
- `github.com/rivo/tview`: Terminal-based interactive viewer.
- `github.com/spf13/pflag`: Flag parsing.
- `github.com/spf13/viper`: Configuration management.

## Contributing

If you'd like to contribute to AMTUI, feel free to submit pull requests or open issues on the [GitHub repository](https://github.com/pehlicd/amtui). Your feedback and contributions are highly appreciated!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Developed by [Furkan Pehlivan](https://github.com/pehlicd) - [Project Repository](https://github.com/pehlicd/amtui)