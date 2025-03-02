package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set the environment variable to use the test configuration file
	os.Setenv("STAGE", "test")

	// Load the configuration
	config := LoadConfig()

	// Assert the configuration values
	assert.Equal(t, "8081", config.Port)
	assert.Equal(t, 5, config.Requests)
	assert.Equal(t, 5*time.Second, config.Duration)
	assert.Equal(t, 5*time.Minute, config.Message.Expiration)
	assert.Equal(t, 5, config.Message.MinLength)
	assert.Equal(t, 500, config.Message.MaxLength)
}
