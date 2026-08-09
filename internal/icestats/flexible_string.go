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
	"errors"
	"strconv"
)

type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case nil:
		*fs = ""
	case string:
		*fs = FlexibleString(value)
	case float64:
		*fs = FlexibleString(strconv.FormatFloat(value, 'f', -1, 64))
	case bool:
		*fs = FlexibleString(strconv.FormatBool(value))
	default:
		return errors.New("invalid flexible string type")
	}
	return nil
}
