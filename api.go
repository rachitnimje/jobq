package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (js *JobServer) EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	var savedJob Job
	err := js.DB.QueryRow(`INSERT INTO jobs(payload, status, created_at) 
		VALUES($1, $2, $3) 
		RETURNING id, payload, status, created_at`, "some dummy work", StatusPending, time.Now()).
		Scan(&savedJob.ID, &savedJob.Payload, &savedJob.Status, &savedJob.CreatedAt)

	if err != nil {
		log.Println("Error creating job:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js.WG.Add(1)
	js.Queue <- &savedJob

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&savedJob)
	if err != nil {
		log.Fatal("Error encoding job to json: ", err)
		return
	}
}

func (js *JobServer) StatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var status string
	err := js.DB.QueryRow(`SELECT status FROM jobs WHERE id = $1`, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			log.Fatal("Job not found:", err)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(status))
	if err != nil {
		log.Fatal("Error writing status: ", err)
		return
	}
}
