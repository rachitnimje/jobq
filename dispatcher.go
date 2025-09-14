package main

import (
	"strconv"
	"time"
)

func StartDispatcher(totalJobs int, server *JobServer) {
	for i := 1; i <= totalJobs; i++ {
		job := &Job{
			ID:        strconv.Itoa(i),
			Payload:   "this is some dummy work",
			Status:    StatusPending,
			CreatedAt: time.Now(),
		}

		server.WG.Add(1)
		server.Queue <- job
		server.Map[job.ID] = job
	}
}
