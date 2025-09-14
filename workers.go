package main

import (
	"fmt"
	"sync"
	"time"
)

func StartWorkerPool(n int, server *JobServer, workersWg *sync.WaitGroup) {
	for id := 1; id <= n; id++ {
		workersWg.Add(1)
		go worker(id, server, workersWg)
	}
}

func worker(id int, server *JobServer, workersWg *sync.WaitGroup) {
	defer workersWg.Done()

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

	fmt.Printf("Worker %d shutting down...\n", id)
}
