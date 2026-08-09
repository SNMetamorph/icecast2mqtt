package config

import (
	"encoding/json"
	"os"
)

type HADiscoveryConfig struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

type TargetConfigEntry struct {
	URL            string             `json:"url"`
	MQTTTopic      string             `json:"mqtt_topic"`
	UpdateInterval TimeDuration       `json:"update_interval"`
	HADiscovery    *HADiscoveryConfig `json:"ha_discovery,omitempty"`
}

type MqttSettings struct {
	BrokerURL           string
	User                string
	Password            string
	ClientID            string `json:"client_id"`
	BaseTopic           string `json:"base_topic"`
	AutoDiscoveryPrefix string `json:"ha_autodiscovery_prefix"`
}

type Config struct {
	MQTT    MqttSettings
	Targets []TargetConfigEntry
}

func Load(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	var object struct {
		MQTTSettings MqttSettings        `json:"mqtt"`
		Targets      []TargetConfigEntry `json:"targets"`
	}
	if err := json.NewDecoder(file).Decode(&object); err != nil {
		return nil, err
	}

	cfg := &Config{
		MQTT: MqttSettings{
			BrokerURL:           os.Getenv("MQTT_BROKER"),
			User:                os.Getenv("MQTT_USERNAME"),
			Password:            os.Getenv("MQTT_PASSWORD"),
			BaseTopic:           object.MQTTSettings.BaseTopic,
			AutoDiscoveryPrefix: object.MQTTSettings.AutoDiscoveryPrefix,
		},
		Targets: object.Targets,
	}
	return cfg, nil
}

func (p *TargetConfigEntry) GetAutoDiscoveryStatus() string {
	if p.HADiscovery != nil {
		return "enabled"
	}
	return "disabled"
}
