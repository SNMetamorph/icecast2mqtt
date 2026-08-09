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

package icestats

import (
	"encoding/json"
	"path"
	"strings"
)

type IcecastSourceMetadata struct {
	Artist      string         `json:"artist"`
	Title       string         `json:"title"`
	Album       string         `json:"album,omitempty"`
	Comment     string         `json:"comment,omitempty"`
	Date        FlexibleString `json:"date,omitempty"`
	Genre       FlexibleString `json:"genre,omitempty"`
	TrackNumber FlexibleString `json:"tracknumber,omitempty"`
}

type IcecastSource struct {
	AudioChannels     int                    `json:"audio_channels"`
	AudioSampleRate   int                    `json:"audio_samplerate"`
	DisplayTitle      FlexibleString         `json:"display-title"`
	Genre             string                 `json:"genre"`
	InstanceUUID      string                 `json:"instance_uuid"`
	ListenersPeak     int                    `json:"listener_peak"`
	Listeners         int                    `json:"listeners"`
	ListenURL         string                 `json:"listenurl"`
	ServerDescription string                 `json:"server_description"`
	ServerName        string                 `json:"server_name"`
	ServerType        string                 `json:"server_type"`
	StreamStartDate   IcecastTime            `json:"stream_start"`
	Subtype           string                 `json:"subtype"`
	ContentType       string                 `json:"content-type"`
	Artist            *FlexibleString        `json:"artist,omitempty"`
	Title             *FlexibleString        `json:"title,omitempty"`
	Metadata          *IcecastSourceMetadata `json:"metadata,omitempty"`
}

// Extract mountpoint from ListenURL ("/stream.mp3" -> "stream")
// TODO: I guess there is the better way to do it. Especially in case of multi-level URL
func (s *IcecastSource) GetMountpoint() string {
	if idx := strings.LastIndex(s.ListenURL, "/"); idx != -1 {
		pathName := s.ListenURL[idx+1:]
		return strings.TrimSuffix(pathName, path.Ext(pathName))
	}
	return strings.TrimSuffix(s.ListenURL, path.Ext(s.ListenURL))
}

func (s *IcecastSource) UnmarshalJSON(data []byte) error {
	type Alias IcecastSource

	aux := &struct {
		StreamStartISO IcecastTime `json:"stream_start_iso8601"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if !aux.StreamStartISO.Time.IsZero() {
		s.StreamStartDate = aux.StreamStartISO
	}

	return nil
}
