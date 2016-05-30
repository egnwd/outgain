package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/egnwd/outgain/server/engine"
	"gopkg.in/antage/eventsource.v1"
)

func UpdatesHandler(engine *engine.Engine) http.Handler {
	events := eventsource.New(nil, nil)
	go func() {
		for event := range engine.Events {
			packet, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("JSON serialization failed %v", err)
			} else {
				events.SendEventMessage(string(packet), event.Type, "")
			}
		}
	}()

	return events
}
