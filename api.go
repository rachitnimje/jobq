package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func (js *JobServer) EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	job := &Job{
		ID:        fmt.Sprintf("%d", len(js.Map)+1),
		Payload:   "dummy payload",
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	js.Mutex.Lock()
	js.Map[job.ID] = job
	js.Mutex.Unlock()

	js.WG.Add(1)
	js.Queue <- job

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(job)
	if err != nil {
		log.Fatal("Error encoding job to json: ", err)
		return
	}
}

func (js *JobServer) StatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	js.Mutex.RLock()
	job, exists := js.Map[id]
	js.Mutex.RUnlock()

	if !exists {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(job)
	if err != nil {
		log.Fatal("Error encoding job to json: ", err)
		return
	}
}
