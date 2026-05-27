package onvifeventer

import (
	"context"
	"log"
	"net/http"
	"time"

	goonvif "github.com/rusl222/onvif"
	"github.com/rusl222/onvif/event"
	sdk "github.com/rusl222/onvif/sdk/event"
)

type MotionEventListener interface {
	MotionDetect(time.Time)
}

type OnvifEventer struct {
	conf     OnvifEventerConfig
	listener MotionEventListener
}

func New(conf OnvifEventerConfig, listener MotionEventListener) *OnvifEventer {
	return &OnvifEventer{
		conf:     conf,
		listener: listener,
	}
}

func (ev *OnvifEventer) Run(ctx context.Context) {
	var dev *goonvif.Device
	var err error
	var eventLink string

	select {
	case <-ctx.Done():
		if eventLink != "" {
			sdk.Call_Unsubscribe(ctx, dev, eventLink, event.Unsubscribe{})
		}
		return
	default:
		// step 1
		for {
			//Getting an camera instance
			dev, err = goonvif.NewDevice(goonvif.DeviceParams{
				Xaddr:      ev.conf.Host,
				Username:   ev.conf.User,
				Password:   ev.conf.Password,
				HttpClient: new(http.Client),
			})
			if err != nil {
				log.Printf("[err] OnvifEventer - %v\t Sleeping a minute", err)
			} else {
				break
			}
			time.Sleep(time.Minute)
		}

		//step 2
		for {
			//Preparing commands
			resp, err := sdk.Call_CreatePullPointSubscription(ctx, dev, event.CreatePullPointSubscription{})

			if err != nil {
				log.Printf("[err] OnvifEventer - %v\t Sleeping a minute", err)
			} else {
				eventLink = string(resp.SubscriptionReference.Address)
				break
			}
			time.Sleep(time.Minute)
		}

		// polling Messages
		for {
			resp2, err := sdk.Call_PullMessages(ctx, dev, eventLink, event.PullMessages{MessageLimit: 1})
			if err != nil {
				log.Printf("[err] OnvifEventer - %v", err)
			}

			for _, mes := range resp2.NotificationMessage {
				t, err := time.Parse(time.RFC3339, string(mes.Message.Message.UtcTime))
				localTime := t.Local()
				if err != nil {
					log.Printf("[err] OnvifEventer - %v", err)
				} else {
					if time.Since(localTime) < ev.conf.PollingDuration {
						message := mes.Message.Message.Data.Name + ":" + string(mes.Message.Message.Data.Value)
						log.Print(message)
						if message == "IsMotion:true" {
							ev.listener.MotionDetect(localTime)
						}

					}
				}
			}
			time.Sleep(ev.conf.PollingDuration)
		}
	}
}
