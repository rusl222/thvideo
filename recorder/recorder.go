package recorder

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Recorder struct {
	conf     RecorderConfig
	listener NewRecordListener
}

type NewRecordListener interface {
	NewRecordReady(filename string, recordDuration time.Duration)
}

func New(conf RecorderConfig, listener NewRecordListener) *Recorder {
	return &Recorder{
		conf:     conf,
		listener: listener,
	}
}

func (rec *Recorder) StartRecord(ctx context.Context) error {

	// Ensure the root directory exists
	if _, err := os.Stat(rec.conf.RootDir); os.IsNotExist(err) {
		if err := os.MkdirAll(rec.conf.RootDir, os.ModePerm); err != nil {
			log.Printf("[err] %v", err)
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():

		default:
			t := time.Now()

			toodayDir := filepath.Join(rec.conf.RootDir, t.Format("2006-01-02"))
			fileName := filepath.Join(toodayDir, fmt.Sprintf("%s.mp4", t.Format("15_04_05")))

			// Create directories if it doesn't exist
			if _, err := os.Stat(toodayDir); os.IsNotExist(err) {
				if err := os.MkdirAll(toodayDir, os.ModePerm); err != nil {
					log.Printf("[err] %v", err)
					return err
				}
			}

			// Start recording
			args := []string{
				"-i", rec.conf.VideoLink,
				"-c", "copy", "-t",
				fmt.Sprintf("%.0f", rec.conf.RecordsDuration.Seconds()),
				"-n", fileName,
			}
			cmd := exec.Command("ffmpeg", args...)

			if err := cmd.Run(); err != nil {
				log.Printf("[err] %v", err)
				time.Sleep(rec.conf.RecordsDuration)
			} else {
				go rec.listener.NewRecordReady(fileName, rec.conf.RecordsDuration)
			}
		}
	}
}
