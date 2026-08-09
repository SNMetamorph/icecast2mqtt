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
