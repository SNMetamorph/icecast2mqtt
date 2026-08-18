/*
 * icecast2mqtt - MQTT bridge for Icecast radio streams, compatible with Home Assistant
 *
 * Copyright (C) 2026 SNMetamorph
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package publisher

import (
	"encoding/json"
	"fmt"
	"icecast2mqtt/internal/buildinfo"
	"icecast2mqtt/internal/config"
	"icecast2mqtt/internal/icestats"
	"log"
	"runtime"
)

var origin = HADeviceOrigin{
	Name:            "icecast2mqtt",
	SoftwareVersion: fmt.Sprintf("%s (%s / %s / %s / %s)", buildinfo.Version, buildinfo.Date, buildinfo.Commit, runtime.GOARCH, runtime.GOOS),
	URL:             "https://github.com/SNMetamorph/icecast2mqtt",
}

type HADiscoveryPayload struct {
	Name              string          `json:"name"`
	StateTopic        string          `json:"state_topic"`
	UniqueTopicID     string          `json:"unique_id"`
	UnitOfMeasurement string          `json:"unit_of_measurement,omitempty"`
	EntityCategory    string          `json:"entity_category,omitempty"`
	DeviceClass       string          `json:"device_class,omitempty"`
	StateClass        string          `json:"state_class,omitempty"`
	Icon              string          `json:"icon,omitempty"`
	Origin            *HADeviceOrigin `json:"origin,omitempty"`
	Device            HADevice        `json:"device"`
}

type HADeviceOrigin struct {
	Name            string `json:"name"`
	SoftwareVersion string `json:"sw"`
	URL             string `json:"url"`
}

type HADevice struct {
	Identifiers     []string `json:"identifiers"`
	Name            string   `json:"name"`
	SoftwareVersion string   `json:"sw_version,omitempty"`
	Manufactuter    string   `json:"manufacturer,omitempty"`
	Model           string   `json:"model,omitempty"`
}

func (p *Publisher) postServerMetrics(baseTopic string, prefix string, cfg config.TargetConfigEntry, device *HADevice) {
	for _, s := range serverMetrics {
		discoveryTopic := fmt.Sprintf("%s/sensor/%s/%s/config", prefix, cfg.HADiscovery.DeviceID, s.SubTopic)
		stateTopic := fmt.Sprintf("%s/%s/%s", baseTopic, cfg.MQTTTopic, s.SubTopic)

		payload := HADiscoveryPayload{
			Name:              fmt.Sprintf("%s %s", cfg.HADiscovery.DeviceName, s.Name),
			StateTopic:        stateTopic,
			UniqueTopicID:     fmt.Sprintf("%s_%s", cfg.HADiscovery.DeviceID, s.SubTopic),
			DeviceClass:       s.DeviceClass,
			StateClass:        s.StateClass,
			EntityCategory:    "diagnostic",
			UnitOfMeasurement: s.Unit,
			Icon:              s.Icon,
			Device:            *device,
			Origin:            &origin,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal HA discovery payload: %v", err)
			continue
		}

		token := p.client.Publish(discoveryTopic, 1, true, data)
		token.Wait()
		if token.Error() != nil {
			log.Printf("[ERROR] Failed to publish HA Discovery to %s: %v", discoveryTopic, token.Error())
		}
	}
}

func (p *Publisher) postStreamMetrics(baseTopic string, prefix string, cfg config.TargetConfigEntry, stream *icestats.IcecastSource, device *HADevice) {
	for _, s := range streamMetrics {
		discoveryTopic := fmt.Sprintf("%s/sensor/%s/%s_%s/config", prefix, cfg.HADiscovery.DeviceID, stream.GetMountpoint(), s.SubTopic)
		stateTopic := fmt.Sprintf("%s/%s/%s/%s", baseTopic, cfg.MQTTTopic, stream.GetMountpoint(), s.SubTopic)

		payload := HADiscoveryPayload{
			Name:              fmt.Sprintf("%s %s", stream.GetMountpoint(), s.Name),
			StateTopic:        stateTopic,
			UniqueTopicID:     fmt.Sprintf("%s_%s_%s", cfg.HADiscovery.DeviceID, stream.GetMountpoint(), s.SubTopic),
			DeviceClass:       s.DeviceClass,
			StateClass:        s.StateClass,
			UnitOfMeasurement: s.Unit,
			Icon:              s.Icon,
			Device:            *device,
			Origin:            &origin,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal HA discovery payload: %v", err)
			continue
		}

		token := p.client.Publish(discoveryTopic, 1, true, data)
		token.Wait()
		if token.Error() != nil {
			log.Printf("[ERROR] Failed to publish HA Discovery to %s: %v", discoveryTopic, token.Error())
		}
	}
}

func (p *Publisher) PostHADeviceDiscovery(baseTopic string, prefix string, cfg config.TargetConfigEntry, stats *icestats.IcestatsObject) {
	device := HADevice{
		Identifiers:     []string{cfg.HADiscovery.DeviceID},
		Name:            cfg.HADiscovery.DeviceName,
		SoftwareVersion: origin.SoftwareVersion,
		Manufactuter:    "SNMetamorph",
		Model:           origin.Name,
	}

	p.postServerMetrics(baseTopic, prefix, cfg, &device)
	for _, source := range stats.Sources {
		p.postStreamMetrics(baseTopic, prefix, cfg, &source, &device)
	}
}
