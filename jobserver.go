package main

import (
	"database/sql"
	"sync"
)

type JobServer struct {
	Queue chan *Job
	Mutex *sync.RWMutex
	WG    *sync.WaitGroup
	DB    *sql.DB
}
