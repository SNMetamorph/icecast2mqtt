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
