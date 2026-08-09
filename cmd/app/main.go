package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"icecast2mqtt/internal/config"
	"icecast2mqtt/internal/fetcher"
	"icecast2mqtt/internal/publisher"
	"icecast2mqtt/internal/worker"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	author  = "unknown"
)

func main() {
	banner()

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	opts := mqtt.NewClientOptions().
		AddBroker(cfg.MQTT.BrokerURL).
		SetUsername(cfg.MQTT.User).
		SetPassword(cfg.MQTT.Password).
		SetClientID(cfg.MQTT.ClientID)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker (%s): %v", cfg.MQTT.BrokerURL, token.Error())
	}
	defer mqttClient.Disconnect(250)
	log.Printf("Connected to MQTT broker: %s", cfg.MQTT.BrokerURL)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	httpFetcher := fetcher.New()
	mqttPublisher := publisher.New(mqttClient)

	for _, target := range cfg.Targets {
		w := worker.New(*cfg, target, httpFetcher, mqttPublisher)
		go w.Start(ctx)
	}

	log.Printf("Started %d polling workers. Press Ctrl+C to stop.", len(cfg.Targets))

	for i, target := range cfg.Targets {
		log.Printf("[%d] %s | sub-topic \"%s\" | update interval %s | MQTT autodiscovery %s",
			i+1,
			target.URL,
			target.MQTTTopic,
			target.UpdateInterval.Duration().String(),
			target.GetAutoDiscoveryStatus(),
		)
	}

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}

func banner() {
	fmt.Printf("\n")
	fmt.Printf("  icecast2mqtt - MQTT bridge for Icecast radio streams, compatible with Home Assistant\n")
	fmt.Printf("  Version      : %s (%s / %s / %s / %s)\n", version, date, commit, runtime.GOARCH, runtime.GOOS)
	fmt.Printf("  Website      : https://github.com/SNMetamorph/icecast2mqtt\n")
	fmt.Printf("\n")
}
