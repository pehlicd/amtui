package pkg

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
