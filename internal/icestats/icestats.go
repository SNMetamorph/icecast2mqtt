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
