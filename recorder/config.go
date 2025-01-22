package recorder

import (
	"encoding/json"
	"time"
)

type RecorderConfig struct {
	RootDir         string
	VideoLink       string
	RecordsDuration time.Duration
}

func (c *RecorderConfig) MarshalJSON() ([]byte, error) {
	type Alias RecorderConfig
	return json.Marshal(&struct {
		RecordsDuration float64 `json:"RecordsDuration"`
		*Alias
	}{
		RecordsDuration: c.RecordsDuration.Seconds(),
		Alias:           (*Alias)(c),
	})
}

func (c *RecorderConfig) UnmarshalJSON(data []byte) error {
	type Alias RecorderConfig
	aux := &struct {
		RecordsDuration float64 `json:"RecordsDuration"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	c.RecordsDuration = time.Duration(aux.RecordsDuration) * time.Second
	return nil
}
