package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	var config Config
	var streams string
	var logInterval int

	flag.StringVar(&streams, "streams", getEnv("RMQ_STREAMS", "stream1"), "Comma-separated list of streams to consume from")
	flag.StringVar(&config.Host, "host", getEnv("RMQ_HOST", "localhost"), "RabbitMQ host")
	flag.StringVar(&config.Vhost, "vhost", getEnv("RMQ_VHOST", "/"), "RabbitMQ vhost")
	flag.IntVar(&config.Port, "port", getEnvAsInt("RMQ_PORT", 5552), "RabbitMQ stream port")
	flag.StringVar(&config.User, "user", getEnv("RMQ_USER", "guest"), "RabbitMQ user")
	flag.StringVar(&config.Password, "password", getEnv("RMQ_PASSWORD", "guest"), "RabbitMQ password")
	flag.IntVar(&config.MaxConsumers, "max-consumers", getEnvAsInt("RMQ_MAX_CONSUMERS", 5), "Max consumers per stream")
	flag.DurationVar(&config.Timeout, "timeout", getEnvAsDuration("RMQ_TIMEOUT", 10*time.Second), "Producer timeout")
	flag.IntVar(&config.InitialCredits, "initial-credits", getEnvAsInt("RMQ_INITIAL_CREDITS", 100), "Initial credits")
	flag.IntVar(&logInterval, "log-interval", getEnvAsInt("LOG_INTERVAL", 5), "Log interval in seconds")
	flag.BoolVar(&config.CRCCheck, "crc-check", getEnvAsBool("RMQ_CRC_CHECK", true), "Enable CRC check for consumers")

	flag.Parse()

	config.Streams = strings.Split(streams, ",")
	config.LogInterval = time.Duration(logInterval) * time.Second

	app := NewApp(config)
	app.Start()
	fmt.Println("Application started. Press Ctrl+C to exit.")

	// Wait for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nShutting down...")
	app.Shutdown()
	fmt.Println("Application stopped.")
}

// Helper functions to read from environment variables
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}
