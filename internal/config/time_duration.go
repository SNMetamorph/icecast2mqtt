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

package config

import (
	"encoding/json"
	"errors"
	"time"
)

type TimeDuration time.Duration

func (td TimeDuration) Duration() time.Duration {
	return time.Duration(td)
}

func (td *TimeDuration) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if value == "" {
		return errors.New("empty time duration string")
	}

	interval, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	*td = TimeDuration(interval.Abs())
	return nil
}
