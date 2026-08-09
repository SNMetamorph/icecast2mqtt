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
	"fmt"
	"icecast2mqtt/internal/icestats"
	"time"
)

type MetricDescription struct {
	SubTopic    string
	Name        string
	DeviceClass string
	StateClass  string
	Unit        string
	Icon        string
}

type ServerMetricPayloadFn func(stats *icestats.IcestatsObject) string
type StreamMetricPayloadFn func(source *icestats.IcecastSource) string

type ServerMetric struct {
	Payload ServerMetricPayloadFn
	MetricDescription
}

type StreamMetric struct {
	Payload StreamMetricPayloadFn
	MetricDescription
}

var serverMetrics []ServerMetric = []ServerMetric{
	{
		Payload: func(stats *icestats.IcestatsObject) string {
			return stats.ServerID
		},
		MetricDescription: MetricDescription{
			SubTopic:    "version",
			Name:        "Version",
			DeviceClass: "",
			StateClass:  "",
			Unit:        "",
			Icon:        "",
		},
	},
	{
		Payload: func(stats *icestats.IcestatsObject) string {
			return stats.StartDate.Format(time.RFC3339)
		},
		MetricDescription: MetricDescription{
			SubTopic:    "server_start_date",
			Name:        "Server Start Date",
			DeviceClass: "timestamp",
			StateClass:  "",
			Unit:        "",
			Icon:        "mdi:calendar-clock",
		},
	},
	{
		Payload: func(stats *icestats.IcestatsObject) string {
			return fmt.Sprintf("%d", len(stats.Sources))
		},
		MetricDescription: MetricDescription{
			SubTopic:    "active_sources",
			Name:        "Sources Count",
			DeviceClass: "",
			StateClass:  "measurement",
			Unit:        "",
			Icon:        "mdi:broadcast",
		},
	},
}

var streamMetrics []StreamMetric = []StreamMetric{
	{
		Payload: func(source *icestats.IcecastSource) string {
			return fmt.Sprintf("%d", source.Listeners)
		},
		MetricDescription: MetricDescription{
			SubTopic:    "listeners",
			Name:        "Listeners",
			DeviceClass: "",
			StateClass:  "measurement",
			Unit:        "",
			Icon:        "mdi:broadcast",
		},
	},
	{
		Payload: func(source *icestats.IcecastSource) string {
			return source.StreamStartDate.Format(time.RFC3339)
		},
		MetricDescription: MetricDescription{
			SubTopic:    "stream_start_date",
			Name:        "Stream Start Date",
			DeviceClass: "timestamp",
			StateClass:  "",
			Unit:        "",
			Icon:        "mdi:calendar-clock",
		},
	},
	{
		Payload: func(source *icestats.IcecastSource) string {
			if source.Artist != nil && source.Title != nil {
				return fmt.Sprintf("%s - %s", *source.Artist, *source.Title)
			}
			return ""
		},
		MetricDescription: MetricDescription{
			SubTopic:    "track",
			Name:        "Current Track",
			DeviceClass: "",
			StateClass:  "",
			Unit:        "",
			Icon:        "mdi:music",
		},
	},
}
