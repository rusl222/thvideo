package onvifeventer

import (
	"encoding/json"
	"time"
)

type OnvifEventerConfig struct {
	Host            string
	User            string
	Password        string
	PollingDuration time.Duration
}

func (c *OnvifEventerConfig) UnmarshalJSON(data []byte) error {
	type Alias OnvifEventerConfig
	aux := &struct {
		PollingDuration float64 `json:"PollingDuration"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	c.PollingDuration = time.Duration(aux.PollingDuration) * time.Second
	return nil
}

func (c *OnvifEventerConfig) MarshalJSON() ([]byte, error) {
	type Alias OnvifEventerConfig
	return json.Marshal(&struct {
		PollingDuration float64 `json:"PollingDuration"`
		*Alias
	}{
		PollingDuration: c.PollingDuration.Seconds(),
		Alias:           (*Alias)(c),
	})
}
