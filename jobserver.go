package main

import "sync"

type JobServer struct {
	Queue chan *Job
	Map   map[string]*Job
	Mutex *sync.RWMutex
	WG    *sync.WaitGroup
}
