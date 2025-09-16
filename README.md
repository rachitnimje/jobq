# Distributed Job Queue – Phase 1 (Concurrent Prototype)

## Overview
This is **Phase 1** of a distributed job queue project, implemented entirely in Go.  
In this phase, the focus is on **concurrency within a single machine**, using Go’s goroutines and channels.

The system allows:
- Jobs to be **enqueued** via an HTTP API.
- A **worker pool** to process jobs concurrently.
- Job statuses (`pending → processing → completed`) to be queried via an API.
- **Graceful shutdown** on `Ctrl+C` (SIGINT/SIGTERM), ensuring all jobs are finished before exit.

This sets the foundation for Phase 2, where Kafka and Postgres will be introduced for true distribution across machines.

---

## Features
- **HTTP API**
    - `POST /enqueue` → Submit a new job.
    - `GET /status/{id}` → Query job status.
- **Worker Pool**
    - Configurable number of workers.
    - Each worker processes jobs concurrently.
- **Job Lifecycle**
    - Tracks `pending`, `processing`, and `completed` states.
- **Graceful Shutdown**
    - Listens for `SIGINT` / `SIGTERM`.
    - Stops HTTP server.
    - Closes job queue.
    - Waits for all workers to finish.
- **Safe Concurrency**
    - `sync.RWMutex` protects shared job map.
    - `sync.WaitGroup` tracks jobs and workers.

---

## Workflow (Phase 1)

1. A client sends `POST /enqueue`.
2. The server creates a new `Job` with status `pending`.
3. The job is stored in a shared `jobsMap` and added to the in-memory `jobQueue` (Go channel).
4. A pool of workers (goroutines) continuously pulls jobs from the queue.
5. Each worker:
    - Updates the job status to `processing`.
    - Simulates work (e.g., `time.Sleep`).
    - Updates the job status to `completed`.
6. Clients can check job progress with `GET /status/{id}`.
7. On shutdown (`Ctrl+C`):
    - The HTTP server stops accepting requests.
    - The job queue is closed.
    - Workers finish remaining jobs, update statuses, and exit cleanly.

---

## Installation & Running

### Prerequisites
- [Go 1.20+](https://go.dev/dl/)

### Clone and Run
```bash
git clone https://github.com/rachitnimje/jobq.git
cd jobq
go run main.go
```

---

## Project Structure
```
├── main.go         # Starts HTTP server, workers, and handles graceful shutdown
├── job.go          # Defines Job struct and statuses
├── workers.go      # Worker pool implementation
├── dispatcher.go   # (Optional) Dispatcher for bulk job injection
├── api.go          # HTTP handlers: /enqueue and /status/{id}
```
