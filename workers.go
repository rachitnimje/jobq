package main

import (
	"fmt"
	"time"
)

func StartWorkerPool(n int, server *JobServer) {
	for id := 1; id <= n; id++ {
		go worker(id, server)
	}
}

func worker(id int, server *JobServer) {
	for job := range server.Queue {
		fmt.Printf("Worker ID %d executing job ID %s.\n", id, job.ID)

		server.Mutex.Lock()
		job.Status = StatusProcessing
		server.Mutex.Unlock()

		time.Sleep(time.Second)

		fmt.Printf("Job ID %s succesfull.\n", job.ID)
		server.Mutex.Lock()
		job.Status = StatusCompleted
		server.Mutex.Unlock()

		server.WG.Done()
	}
}
