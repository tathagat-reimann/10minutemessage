package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Requests int
	Duration time.Duration
	Message  MessageConfig
}

type MessageConfig struct {
	Expiration time.Duration
	MinLength  int
	MaxLength  int
}

func LoadConfig() Config {
	// Read the environment variable stage
	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = ""
	} else {
		stage = "-" + stage
	}

	viper.SetConfigName("config" + stage)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Get configuration values
	port := viper.GetString("server.port")
	requests := viper.GetInt("rate_limit.requests")
	duration := viper.GetDuration("rate_limit.duration")

	return Config{
		Port:     port,
		Requests: requests,
		Duration: duration,
		Message: MessageConfig{
			Expiration: viper.GetDuration("message.expiration"),
			MinLength:  viper.GetInt("message.min_length"),
			MaxLength:  viper.GetInt("message.max_length"),
		},
	}
}
