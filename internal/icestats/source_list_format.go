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
