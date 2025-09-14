package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	server := &JobServer{
		Queue: make(chan *Job, 1000),
		Map:   make(map[string]*Job),
		Mutex: &sync.RWMutex{},
		WG:    &wg,
	}

	// define routes
	router := mux.NewRouter()
	router.HandleFunc("/enqueue", server.EnqueueHandler).Methods("POST")
	router.HandleFunc("/status/{id}", server.StatusHandler).Methods("GET")

	// start the workers
	totalWorkers := 100
	StartWorkerPool(totalWorkers, server)

	// start the http server
	fmt.Println("Server running at :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting http server on port 8080. ", err)
		return
	}
}
