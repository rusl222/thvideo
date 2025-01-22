package filemanager

import (
	"log"
	"os"
	"path"
	"testing"
	"time"
)

func createFile(filename string) string {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create test record file: %v", err)
	}
	file.Close()
	return filename
}

func TestFileManager(t *testing.T) {

	rootDir := "d:/tmp/test_filemanager"
	conf := FileManagerConfig{RootDir: rootDir}
	fm := New(conf)

	// Test NewRecordReady
	day := time.Now()
	camDir := path.Join(rootDir, day.Format("2006-01-02"))

	// Create test directory
	if err := os.MkdirAll(camDir, os.ModePerm); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test record files
	recordFile1 := createFile(path.Join(camDir, "test_record1.mp4"))
	recordFile2 := createFile(path.Join(camDir, "test_record2.mp4"))
	recordFile3 := createFile(path.Join(camDir, "test_record3.mp4"))

	recordsDuration := 2 * time.Second
	fm.NewRecordReady(recordFile1, recordsDuration)
	if fm.prevFile != recordFile1 {
		t.Errorf("Expected to contain %s, but got %s", recordFile1, fm.prevFile)
	}

	time.Sleep(time.Second)

	// Test MotionDetect
	motionTime := time.Now()
	fm.MotionDetect(motionTime)
	if fm.lastMotionTime != motionTime {
		t.Errorf("Expected lastMotionTime to be %v, but got %v", motionTime, fm.lastMotionTime)
	}

	time.Sleep(6 * time.Second)
	fm.NewRecordReady(recordFile2, recordsDuration)
	time.Sleep(time.Second)

	// Test that file not deleteted if lastMotionTime is not older than recordsDuration
	if _, err := os.Stat(recordFile1); os.IsNotExist(err) {
		t.Errorf("Expected record file %s NOT to be deleted, but it deleted", recordFile1)
	}

	time.Sleep(8 * time.Second)
	fm.NewRecordReady(recordFile3, recordsDuration)
	time.Sleep(time.Second)

	// Test that file is deleted if lastMotionTime is older than recordsDuration +5 seconds
	if _, err := os.Stat(recordFile2); !os.IsNotExist(err) {
		t.Errorf("Expected record file %s to be deleted, but it still exists", recordFile2)
	}

	// Test GetFiles
	if err := os.MkdirAll(camDir, os.ModePerm); err != nil {
		t.Fatalf("Failed to create camera directory: %v", err)
	}

	files := fm.GetFiles(day)
	notOk := len(files) != 2 || files[0] != "/test_record1.mp4" || files[1] != "/test_record3.mp4"
	if notOk {
		t.Errorf("Expected files to contain %s and %s, but got %v", "/test_record1.mp4", "/test_record3.mp4", files)
	}

	os.RemoveAll(rootDir) // Clean up
}
