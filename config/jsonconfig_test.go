package config

import (
	"atvideo/filemanager"
	"atvideo/recorder"
	"atvideo/videoplayer"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestJsonConfig_LoadConfig(t *testing.T) {
	// Setup
	configFilePath := filepath.Join(os.TempDir(), "test_config.json")
	configData := `{
		"RootDir": "/tmp",
		"Cams": [
			{
				"Name": "Cam1",
				"FileManager": {},
				"Recorder": {}
			}
		],
		"VideoServer": {}
	}`
	if err := os.WriteFile(configFilePath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	defer os.Remove(configFilePath)

	jsonConfig := NewJsonConfig()

	// Test LoadConfig
	config, err := jsonConfig.LoadConfig(configFilePath)
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	// Validate loaded config
	if config.RootDir != "/tmp" {
		t.Errorf("Expected RootDir to be /tmp, but got %s", config.RootDir)
	}
	if len(config.Cams) != 1 || config.Cams[0].Name != "Cam1" {
		t.Errorf("Expected one Cam with Name Cam1, but got %v", config.Cams)
	}
	if config.Cams[0].FileManager.RootDir != "/tmp/Cam1" {
		t.Errorf("Expected FileManager RootDir to be /tmp/Cam1, but got %s", config.Cams[0].FileManager.RootDir)
	}
	if config.Cams[0].Recorder.RootDir != "/tmp/Cam1" {
		t.Errorf("Expected Recorder RootDir to be /tmp/Cam1, but got %s", config.Cams[0].Recorder.RootDir)
	}
	if config.VideoServer.RootDir != "/tmp" {
		t.Errorf("Expected VideoServer RootDir to be /tmp, but got %s", config.VideoServer.RootDir)
	}
}

func TestJsonConfig_SaveConfig(t *testing.T) {
	// Setup
	configFilePath := filepath.Join(os.TempDir(), "test_config.json")
	defer os.Remove(configFilePath)

	config := &ThVideoConfig{
		RootDir: "/tmp",
		Cams: []CamConfig{
			{
				Name: "Cam1",
				FileManager: filemanager.FileManagerConfig{
					RootDir: "/tmp/Cam1",
				},
				Recorder: recorder.RecorderConfig{
					RootDir: "/tmp/Cam1",
				},
			},
		},
		VideoServer: videoplayer.VideoPlayerBackendConfig{
			RootDir: "/tmp",
		},
	}

	jsonConfig := NewJsonConfig()

	// Test SaveConfig
	if err := jsonConfig.SaveConfig(configFilePath, config); err != nil {
		t.Fatalf("SaveConfig returned an error: %v", err)
	}

	// Validate saved config
	savedConfigData, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("Failed to read saved config file: %v", err)
	}

	var savedConfig ThVideoConfig
	if err := json.Unmarshal(savedConfigData, &savedConfig); err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if savedConfig.RootDir != "/tmp" {
		t.Errorf("Expected RootDir to be /tmp, but got %s", savedConfig.RootDir)
	}
	if len(savedConfig.Cams) != 1 || savedConfig.Cams[0].Name != "Cam1" {
		t.Errorf("Expected one Cam with Name Cam1, but got %v", savedConfig.Cams)
	}

}
