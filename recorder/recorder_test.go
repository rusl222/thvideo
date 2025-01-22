package recorder

import (
	"context"
	"os"
	"testing"
	"time"
)

type MockListener struct {
	records []string
}

func (ml *MockListener) NewRecordReady(filename string, recordDuration time.Duration) {
	ml.records = append(ml.records, filename)
}

func TestStartRecord(t *testing.T) {
	// Setup
	rootDir := "d:\\tmp\\test_records"
	videoLink := "rtsp://watcher:Watch1234@192.168.1.64:554/Streaming/Channels/101"
	recordsDuration := 10 * time.Second

	conf := RecorderConfig{
		RootDir:         rootDir,
		VideoLink:       videoLink,
		RecordsDuration: recordsDuration,
	}

	listener := &MockListener{}
	recorder := New(conf, listener)

	// Create test directory
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(rootDir) // Clean up

	// Start recording in a separate goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := recorder.StartRecord(ctx); err != nil {
			t.Errorf("StartRecord returned an error: %v", err)
		}
	}()

	// Allow some time for recording to start
	time.Sleep(25 * time.Second)
	cancel() // Stop recording

	// Check if a new record was created
	if len(listener.records) == 0 {
		t.Errorf("Expected at least one record to be created, but got none")
	}

	// Check if the file exists
	for _, record := range listener.records {
		t.Logf("Checking record file: %s", record)
		if _, err := os.Stat(record); os.IsNotExist(err) {
			t.Errorf("Expected record file %s to exist, but it does not", record)
		}
	}
}
