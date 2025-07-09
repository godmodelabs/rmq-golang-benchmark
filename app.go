package main

import (
	"fmt"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type App struct {
	config    Config
	consumers map[string]*Consumer
	ticker    *time.Ticker
	done      chan bool
}

type Config struct {
	Streams        []string
	Host           string
	Vhost          string
	Port           int
	User           string
	Password       string
	MaxConsumers   int
	Timeout        time.Duration
	InitialCredits int
	LogInterval    time.Duration
	CRCCheck       bool
}

func NewApp(c Config) *App {
	return &App{
		config:    c,
		consumers: make(map[string]*Consumer),
		done:      make(chan bool),
	}
}

func (a *App) Start() {
	go func() {
		for _, streamName := range a.config.Streams {
			env, err := stream.NewEnvironment(
				stream.NewEnvironmentOptions().
					SetHost(a.config.Host).
					SetVHost(a.config.Vhost).
					SetPort(a.config.Port).
					SetUser(a.config.User).
					SetPassword(a.config.Password).
					SetMaxConsumersPerClient(a.config.MaxConsumers).
					SetRPCTimeout(a.config.Timeout),
			)
			if err != nil {
				fmt.Println(err)
			}

			consumerOptions := stream.NewConsumerOptions().
				SetOffset(stream.OffsetSpecification{}.First()).
				SetInitialCredits(int16(a.config.InitialCredits)).
				SetCRCCheck(a.config.CRCCheck)

			consumer, err := NewConsumer(env, streamName, consumerOptions)
			if err != nil {
				fmt.Printf("error for stream %v: %v\n", streamName, err)
				continue
			}

			fmt.Printf("started consumer for stream %v\n", streamName)
			a.consumers[streamName] = consumer
		}
	}()

	a.ticker = time.NewTicker(a.config.LogInterval)
	go a.logStats()
}

func (a *App) logStats() {
	var lastTotal, lastEntries uint64
	lastLoggedTime := time.Now()

	for {
		select {
		case <-a.done:
			return
		case <-a.ticker.C:
			var totalConsumed, totalEntries uint64
			for _, c := range a.consumers {
				totalConsumed += c.consumed.Load()
				totalEntries += c.totalEntries.Load()
			}

			now := time.Now()
			duration := now.Sub(lastLoggedTime)
			messagesSinceLast := totalConsumed - lastTotal
			entriesSinceLast := totalEntries - lastEntries
			lastTotal = totalConsumed
			lastEntries = totalEntries
			lastLoggedTime = now

			mps := float64(messagesSinceLast) / duration.Seconds()
			avgEntries := float64(entriesSinceLast) / float64(messagesSinceLast)

			if messagesSinceLast == 0 {
				avgEntries = 0
			}

			fmt.Printf("Total consumed: %d, Rate: %.2f msg/s, Avg entries: %.2f\n", totalConsumed, mps, avgEntries)
		}
	}
}

func (a *App) Shutdown() {
	a.ticker.Stop()
	a.done <- true

	var err error
	for name, a := range a.consumers {
		err = a.c.Close()
		if err != nil {
			fmt.Printf("error closing stream %v: %v\n", name, err)
		}
	}

	fmt.Println("stopped rmq consumers")
}
