package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var jobsWg sync.WaitGroup
	var workersWg sync.WaitGroup

	server := &JobServer{
		Queue: make(chan *Job, 1000),
		Map:   make(map[string]*Job),
		Mutex: &sync.RWMutex{},
		WG:    &jobsWg,
	}

	// start the workers
	totalWorkers := 10
	StartWorkerPool(totalWorkers, server, &workersWg)

	// stop channel will receive the signals like ctrl + c or stop server
	stop := make(chan os.Signal, 1)
	// if there is any interrupt(from the ones defined below), os will send that signal to stop channel
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// define routes
	router := mux.NewRouter()
	router.HandleFunc("/enqueue", server.EnqueueHandler).Methods("POST")
	router.HandleFunc("/status/{id}", server.StatusHandler).Methods("GET")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// start the http server in a separate goroutine
	go func() {
		fmt.Println("http server running at :8080")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on :8080. %v\n", err)
			return
		}
	}()

	// main goroutine will stop here to check for stop signal, once received the goroutine will continue
	<-stop

	// create a timeout ctx of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown the http server
	fmt.Println("Shutting down the http server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down the http server: ", err)
	}
	fmt.Println("http server shutdown complete...")

	// close the job queue for new jobs
	close(server.Queue)

	// wait till all the processing jobs are completed
	fmt.Println("completing the remaining jobs...")
	server.WG.Wait()

	// wait till all the workers exits
	workersWg.Wait()

	fmt.Println("Shutdown complete...")
}
