package worker

import (
	"context"
	"fmt"
	"icecast2mqtt/internal/config"
	"icecast2mqtt/internal/fetcher"
	"icecast2mqtt/internal/icestats"
	"icecast2mqtt/internal/publisher"
	"log"
	"time"
)

type Worker struct {
	cfg       config.Config
	target    config.TargetConfigEntry
	fetcher   *fetcher.Fetcher
	publisher *publisher.Publisher
}

func New(cfg config.Config, target config.TargetConfigEntry, f *fetcher.Fetcher, pub *publisher.Publisher) *Worker {
	return &Worker{
		cfg:       cfg,
		target:    target,
		fetcher:   f,
		publisher: pub,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.target.UpdateInterval.Duration())
	defer ticker.Stop()

	w.process()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[WORKER] Stopped worker for topic: %s", w.target.MQTTTopic)
			return
		case <-ticker.C:
			w.process()
		}
	}
}

func (w *Worker) process() {
	rawData, err := w.fetcher.Fetch(w.target.URL)
	if err != nil {
		log.Printf("[ERROR] [%s] Fetch failed: %v", w.target.MQTTTopic, err)
		return
	}

	stats, err := icestats.Parse(rawData)
	if err != nil {
		log.Printf("[ERROR] [%s] Parse failed: %v", w.target.MQTTTopic, err)
		return
	}

	if w.target.HADiscovery != nil {
		w.publisher.PostHADeviceDiscovery(w.cfg.MQTT.BaseTopic, w.cfg.MQTT.AutoDiscoveryPrefix, w.target, &stats)
	}

	instanceTopic := fmt.Sprintf("%s/%s", w.cfg.MQTT.BaseTopic, w.target.MQTTTopic)
	w.publisher.PostInstanceMetrics(instanceTopic, stats)

	for _, source := range stats.Sources {
		sourceTopic := fmt.Sprintf("%s/%s", instanceTopic, source.GetMountpoint())
		w.publisher.PostStreamMetrics(sourceTopic, source)
	}
}
