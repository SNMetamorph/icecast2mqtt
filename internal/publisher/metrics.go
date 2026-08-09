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
