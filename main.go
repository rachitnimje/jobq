package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	server := &JobServer{
		Queue: make(chan *Job, 1000),
		Map:   make(map[string]*Job),
		Mutex: &sync.RWMutex{},
		WG:    &wg,
	}

	totalWorkers := 100
	StartWorkerPool(totalWorkers, server)

	totalJobs := 1000
	startTime := time.Now()

	StartDispatcher(totalJobs, server)
	close(server.Queue)

	wg.Wait()

	elapsed := time.Since(startTime).Seconds()
	throughput := float64(totalJobs) / elapsed
	jobsPerWorker := float64(totalJobs) / float64(totalWorkers)

	fmt.Printf("completed %d jobs by %d workers in %f seconds (%f jobs/sec)\n",
		totalJobs, totalWorkers, elapsed, throughput)
	fmt.Printf("avg jobs per worker: %f\n", jobsPerWorker)

}
