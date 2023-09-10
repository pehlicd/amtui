package pkg

import (
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
  -i, --insecure        For insecurely connecting to Alertmanager
  -v, --version         Show version
  -h, --help            Help
`
)

var (
	versionString string
	fl            = flag.NewFlagSet("amtui", flag.ExitOnError)
	host          = fl.String("host", "localhost", "Alertmanager host")
	port          = fl.StringP("port", "p", "9093", "Alertmanager port")
	insecure      = fl.BoolP("insecure", "i", true, "For insecurely connecting to Alertmanager")
	help          = fl.BoolP("help", "h", false, "Show help")
	version       = fl.BoolP("version", "v", false, "Show version")
	scheme        = "http"
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
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Insecure bool   `yaml:"insecure"`
	Scheme   string `yaml:"scheme"`
}

func initConfig() Config {
	if err := fl.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Initialize Viper
	viper.SetConfigName(".amtui")          // Configuration file name without extension
	viper.SetConfigType("yaml")            // Configuration file type
	viper.AddConfigPath(os.Getenv("HOME")) // Search for the configuration file in the $HOME directory

	if *insecure {
		scheme = "http"
	} else {
		scheme = "https"
	}

	config := Config{
		Host:     *host,
		Port:     *port,
		Insecure: *insecure,
		Scheme:   scheme,
	}

	if *help {
		printHelp(os.Stderr)
	}

	if *version {
		fmt.Printf("amtui version: v%s\n", versionString)
		os.Exit(0)
	}

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
