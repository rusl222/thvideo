package config

import (
	"encoding/json"
	"os"
	"path"
)

type JsonConfig struct {
}

func NewJsonConfig() *JsonConfig {
	return &JsonConfig{}
}

func (*JsonConfig) LoadConfig(filePath string) (*ThVideoConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config ThVideoConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (*JsonConfig) SaveConfig(filePath string, config *ThVideoConfig) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Add indentation for pretty printing
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

// UnmarshalJSON is a custom unmarshaler for ThVideoConfig
// that sets the RootDir for all the Cams and the VideoServer.
func (c *ThVideoConfig) UnmarshalJSON(data []byte) error {
	type Alias ThVideoConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for i := range c.Cams {
		c.Cams[i].FileManager.RootDir = path.Join(c.RootDir, c.Cams[i].Name)
		c.Cams[i].Recorder.RootDir = path.Join(c.RootDir, c.Cams[i].Name)
	}
	c.VideoServer.RootDir = c.RootDir

	return nil
}
