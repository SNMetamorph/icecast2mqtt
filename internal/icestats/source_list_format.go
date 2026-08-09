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

import "encoding/json"

type IcecastSourceList []IcecastSource

func (bl *IcecastSourceList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if data[0] == '[' {
		var list []IcecastSource
		if err := json.Unmarshal(data, &list); err != nil {
			return err
		}
		*bl = list
		return nil
	}

	var single IcecastSource
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	*bl = IcecastSourceList{single}
	return nil
}
