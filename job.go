package main

import "time"

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID        string    `json:"id"`
	Payload   string    `json:"payload"`
	Status    JobStatus `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
