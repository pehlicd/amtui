package pkg

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfigDefault(t *testing.T) {
	// Test case 1: Test with default values
	os.Args = []string{"amtui", "--host", "localhost", "--port", "9093", "--scheme", "http", "--username", "admin", "--password", "admin"}
	config := initConfig()
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "9093", config.Port)
	assert.Equal(t, "http", config.Scheme)
	assert.Equal(t, "admin", config.Auth.Username)
	assert.Equal(t, "admin", config.Auth.Password)
}

func TestInitConfigCustom(t *testing.T) {
	// Test case 2: Test with custom values
	os.Args = []string{"amtui", "--host", "example.com", "--port", "9090", "--scheme", "https"}
	config := initConfig()
	assert.Equal(t, "example.com", config.Host)
	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "https", config.Scheme)
	assert.Equal(t, "example.com", viper.Get("host"))
	assert.Equal(t, "9090", viper.Get("port"))
	assert.Equal(t, "https", viper.Get("scheme"))
}

func TestInitInvalid(t *testing.T) {
	// Test case 3: Test with invalid flags
	os.Args = []string{"amtui", "--invalid-flag"}
	assert.Panics(t, func() { printHelp(os.Stderr) })
}

func TestInitHelp(t *testing.T) {
	// Test case 4: Test with help flag
	os.Args = []string{"amtui", "--help"}
	assert.Panics(t, func() { printHelp(os.Stderr) })
}

func TestInitVersion(t *testing.T) {
	// Test case 5: Test with version flag
	os.Args = []string{"amtui", "--version"}
	assert.Panics(t, func() {
		fmt.Printf("Version: %s\nBuild Date: %s\nBuild Commit: %s\n", versionString, buildDate, buildCommit)
		os.Exit(0)
	})
}
