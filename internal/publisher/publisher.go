package publisher

import (
	"fmt"
	"icecast2mqtt/internal/icestats"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	client mqtt.Client
}

func New(client mqtt.Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) publishString(topic, payload string) {
	token := p.client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Printf("[ERROR] Failed to publish to %s: %v", topic, token.Error())
	}
}

func (p *Publisher) PostInstanceMetrics(baseTopic string, stats icestats.IcestatsObject) {
	for _, m := range serverMetrics {
		p.publishString(fmt.Sprintf("%s/%s", baseTopic, m.SubTopic), m.Payload(&stats))
	}
}

func (p *Publisher) PostStreamMetrics(baseTopic string, source icestats.IcecastSource) {
	for _, m := range streamMetrics {
		payload := m.Payload(&source)
		if len(payload) > 0 {
			p.publishString(fmt.Sprintf("%s/%s", baseTopic, m.SubTopic), payload)
		}
	}
}
