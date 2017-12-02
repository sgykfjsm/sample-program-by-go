package main

import (
	"fmt"
	"time"
)

type Worker struct {
	ID          int
	WorkChan    chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan chan WorkRequest) *Worker {
	return &Worker{
		ID:          id,
		WorkChan:    make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.WorkChan

			select {
			case work := <-w.WorkChan:
				fmt.Printf("worked%d: Received work request, delaying for %f seconds\n", w.ID, work.Delay.Seconds())
				time.Sleep(work.Delay)
				fmt.Printf("worked%d: Hello, %s!\n", w.ID, work.Name)
			case <-w.QuitChan:
				fmt.Printf("worker%d: stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
