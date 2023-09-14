package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfigDefault(t *testing.T) {
	// Test case 1: Test with default values
	os.Args = []string{"amtui", "--host", "localhost", "--port", "9093", "--scheme", "http"}
	config := initConfig()
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "9093", config.Port)
	assert.Equal(t, "http", config.Scheme)
}

func TestInitConfigCustom(t *testing.T) {
	// Test case 2: Test with custom values
	os.Args = []string{"amtui", "--host", "example.com", "--port", "9090", "--scheme", "https"}
	config := initConfig()
	assert.Equal(t, "example.com", config.Host)
	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "https", config.Scheme)
}

func TestInitInvalid(t *testing.T) {
	// Test case 3: Test with invalid flags
	os.Args = []string{"amtui", "--invalid-flag"}
	assert.Panics(t, func() { printHelp(os.Stderr) })
}
