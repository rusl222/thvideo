package main

import (
	"atvideo/filemanager"
	"atvideo/onvifeventer"
	"atvideo/recorder"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"atvideo/config"
)

func main() {

	pathConf := "./config.json"

	// Check if a custom config path was provided
	if len(os.Args) > 1 {
		pathConf = os.Args[1]
	}

	// Load the configuration
	conf, err := config.NewJsonConfig().LoadConfig(pathConf)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Listen for OS interrupt signal to gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Received interrupt signal, shutting down...")
		cancel()
	}()

	c := conf.Cams[0]
	// Start the file manager
	fm := filemanager.New(c.FileManager)

	// Start the onvif eventer
	onv := onvifeventer.New(c.OnvifEventer, fm)
	go onv.Run(ctx)

	// Start the recorder
	rec := recorder.New(c.Recorder, fm)
	rec.StartRecord(ctx)

	// Start the video player backend
	// vp := videoplayer.New(conf.VideoServer, fm)
	// vp.Run(":8080")

}
