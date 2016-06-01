package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/egnwd/outgain/server/engine"
	"github.com/egnwd/outgain/server/lobby"
	"github.com/gorilla/mux"
	"gopkg.in/antage/eventsource.v1"
)

func UpdatesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		l, ok := lobby.GetLobby(uint64(id))
		if !ok {
			log.Println("Updates: Lobby doesn't exist")
			http.Error(w, "Lobby doesn't exist", http.StatusInternalServerError)
			return
		}

		eng := l.Engine.(*engine.Engine)
		events := eventsource.New(nil, nil)

		go func() {
			for event := range eng.Events {
				packet, err := json.Marshal(event.Data)
				if err != nil {
					log.Printf("JSON serialization failed %v", err)
				} else {
					events.SendEventMessage(string(packet), event.Type, "")
					if event.Type == "shutdown" {
						eng.Shutdown()
						lobby.DestroyLobby(l)
					}
				}
			}
		}()

		events.ServeHTTP(w, r)
	})
}
