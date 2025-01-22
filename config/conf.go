package config

import (
	"atvideo/filemanager"
	"atvideo/onvifeventer"
	"atvideo/recorder"
	"atvideo/videoplayer"
)

type CamConfig struct {
	Name         string
	FileManager  filemanager.FileManagerConfig
	OnvifEventer onvifeventer.OnvifEventerConfig
	Recorder     recorder.RecorderConfig
}

type ThVideoConfig struct {
	RootDir     string
	Cams        []CamConfig
	VideoServer videoplayer.VideoPlayerBackendConfig
}
