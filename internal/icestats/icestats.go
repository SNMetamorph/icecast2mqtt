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
	"fmt"
)

type IcestatsObject struct {
	Admin        string            `json:"admin"`
	Host         string            `json:"host"`
	InstanceUUID string            `json:"instance_uuid"`
	Location     string            `json:"location"`
	ServerID     string            `json:"server_id"`
	StartDate    IcecastTime       `json:"server_start"`
	Sources      IcecastSourceList `json:"source"`
}

func (p *IcestatsObject) UnmarshalJSON(data []byte) error {
	type Alias IcestatsObject

	aux := &struct {
		StreamStartISO IcecastTime `json:"server_start_iso8601"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if !aux.StreamStartISO.Time.IsZero() {
		p.StartDate = aux.StreamStartISO
	}

	return nil
}

func Parse(rawJSON []byte) (IcestatsObject, error) {
	var object struct {
		Icestats IcestatsObject `json:"icestats"`
	}

	if err := json.Unmarshal(rawJSON, &object); err != nil {
		return IcestatsObject{}, fmt.Errorf("json unmarshal failed: %w", err)
	}

	return object.Icestats, nil
}
