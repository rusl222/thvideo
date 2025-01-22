package filemanager

import (
	"log"
	"os"
	"time"
)

type FileManagerConfig struct {
	RootDir string
}

// FileManager designed to manage record files.
type FileManager struct {
	conf           FileManagerConfig
	prevFile       string
	lastMotionTime time.Time
}

func New(conf FileManagerConfig) *FileManager {
	return &FileManager{
		conf:     conf,
		prevFile: "",
	}
}

// MotionDetected - method for receiving a Motion Detection Signal.
func (fm *FileManager) MotionDetect(motionTime time.Time) {
	log.Printf("new motion detected %v", motionTime)
	fm.lastMotionTime = motionTime
}

// NewRecordReady - method for receiving filename of new record.
// When a new file is received, it checks the time of the last MotionEvent,
// if it is older than the duration of the recordfile, it deletes the previous one.
func (fm *FileManager) NewRecordReady(filename string, recordsDuration time.Duration) {
	log.Printf("new record ready %s %v", filename, recordsDuration)
	if time.Since(fm.lastMotionTime) > recordsDuration+5*time.Second {
		fm.deleteLast()
	}
	fm.prevFile = filename
}

func (fm *FileManager) deleteLast() error {
	err := os.Remove(fm.prevFile)
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("file <%s> - deleted", fm.prevFile)
	}
	return err
}

// GetFiles - method for getting list of files.
func (fm *FileManager) GetFiles(date time.Time) []string {
	ret := []string{}
	camDir := fm.conf.RootDir + "/" + date.Format("/2006-01-02")
	entries, err := os.ReadDir(camDir)
	if err != nil {
		log.Printf("[err] %v", err)
		return ret
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ret = append(ret, "/"+entry.Name())
		}
	}
	return ret
}
