package pkg

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	helpMessage = `Usage: amtui [options]

Options:
  --host                Alertmanager host
  -p, --port            Alertmanager port
  -s, --scheme          Alertmanager scheme (http or https)
  -v, --version         Show version
  -h, --help            Help
`
)

var (
	versionString, buildDate, buildCommit string
	fl                                    = flag.NewFlagSet("amtui", flag.ExitOnError)
	host                                  = fl.String("host", "", "Alertmanager host")
	port                                  = fl.StringP("port", "p", "", "Alertmanager port")
	scheme                                = fl.StringP("scheme", "s", "", "Alertmanager scheme (http or https)")
	username                              = fl.StringP("username", "", "", "Alertmanager username for basic auth")
	password                              = fl.StringP("password", "", "", "Alertmanager password for basic auth")
	help                                  = fl.BoolP("help", "h", false, "Show help")
	version                               = fl.BoolP("version", "v", false, "Show version")
)

func printHelp(w io.Writer) {
	_, err := fmt.Fprint(w,
		helpMessage+"\n")
	if err != nil {
		log.Fatalf("Error writing help to stdout: %v", err)
	}
	os.Exit(0)
}

type Config struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Scheme string `yaml:"scheme"`
	Auth   Auth   `yaml:"auth"`
}

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func initConfig() Config {
	if err := fl.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Initialize Viper
	viper.SetConfigName(".amtui")          // Configuration file name without extension
	viper.SetConfigType("yaml")            // Configuration file type
	viper.AddConfigPath(os.Getenv("HOME")) // Search for the configuration file in the $HOME directory

	// Print help and exit
	if *help {
		printHelp(os.Stderr)
	}

	// Print version and exit
	if *version {
		fmt.Printf("Version: %s\nBuild Date: %s\nBuild Commit: %s\n", versionString, buildDate, buildCommit)
		os.Exit(0)
	}

	// Scheme must be http or https
	if *scheme != "https" && *scheme != "http" && *scheme != "" {
		log.Fatalf("Error: scheme must be http or https. Got: %s\n", *scheme)
	}

	var config Config
	if *username == "" && *password == "" {
		config = Config{
			Host:   *host,
			Port:   *port,
			Scheme: *scheme,
		}
	} else {
		config = Config{
			Host:   *host,
			Port:   *port,
			Scheme: *scheme,
			Auth: Auth{
				Username: *username,
				Password: *password,
			},
		}
	}

	// if flags are set, overwrite config file
	if config.Host != "" && config.Port != "" && config.Scheme != "" {
		viper.Set("host", config.Host)
		viper.Set("port", config.Port)
		viper.Set("scheme", config.Scheme)
		if config.Auth.Username != "" {
			viper.Set("auth.username", config.Auth.Username)
			viper.Set("auth.password", config.Auth.Password)
		}
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// Handle errors when the configuration file is not found or is invalid
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("Config file not found, using defaults.")
			// Write the default configuration to a new file
			if err := viper.SafeWriteConfig(); err != nil {
				log.Fatalf("Error creating config file: %v", err)
			}
		}
	}

	// Merge flags into the configuration
	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		log.Fatalf("Error binding flags: %v", err)
	}

	// Unmarshal the configuration into Config struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return config
}
