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
