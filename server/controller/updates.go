package controller

import (
	"encoding/json"
	"github.com/egnwd/outgain/server/engine"
	"gopkg.in/antage/eventsource.v1"
	"log"
	"net/http"
)

func UpdatesHandler(engine *engine.Engine) http.Handler {
	events := eventsource.New(nil, nil)
	go func() {
		for update := range engine.Updates {
			packet, err := json.Marshal(update)
			if err != nil {
				log.Printf("JSON serialization failed %v", err)
			} else {
				events.SendEventMessage(string(packet), "", "")
			}
		}
	}()

	return events
}
