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
	"fmt"
	"strings"
	"time"
)

const icecastAlternativeFormat = "02/Jan/2006:15:04:05 -0700"
const iso8601Format = "2006-01-02T15:04:05-0700"

type IcecastTime struct {
	time.Time
}

func (ft *IcecastTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}

	// used in modern versions of Icecast
	t, err := time.Parse(iso8601Format, s)
	if err == nil {
		ft.Time = t
		return nil
	}

	// used in older versions of vanilla Icecast
	t, err = time.Parse(time.RFC1123Z, s)
	if err == nil {
		ft.Time = t
		return nil
	}

	// used in Icecast-kh
	t, err = time.Parse(icecastAlternativeFormat, s)
	if err == nil {
		ft.Time = t
		return nil
	}

	return fmt.Errorf("Failed to parse time %q", s)
}

func (ft IcecastTime) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%s"`, ft.Time.Format(time.RFC3339)), nil
}
