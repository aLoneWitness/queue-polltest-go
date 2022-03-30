package main

import (
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"net/http"
	"sync"
	"time"
)

func main() {
	sessionMap := sync.Map{}

	cJob := cron.New()
	cJob.AddFunc("@every 5s", func() {
		now := time.Now().Unix()

		sessionMap.Range(func(key, value interface{}) bool {
			if value.(time.Time).Unix() < now {
				sessionMap.Delete(key)
			}
			return true
		})
	})

	cJob.Start()

	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("Healthy"))
	})

	http.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}
		uniqueId := r.Header.Get("LB_HEADER_AFFINITY")
		if len(uniqueId) == 0 || !IsValidUUID(uniqueId) {
			w.WriteHeader(400)
			w.Write([]byte("Header LB_HEADER_AFFINITY not valid: " + uniqueId))
			return
		}

		sessionMap.Store(uniqueId, time.Now().Add(5*time.Second))
	})

	http.ListenAndServe(":8080", nil)
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
