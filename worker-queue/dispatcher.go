package main

import "fmt"

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworker int) {
	WorkerQueue = make(chan chan WorkRequest, nworker)

	for i := 0; i < nworker; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received work request")
				go func() {
					workerChan := <-WorkerQueue
					fmt.Println("Dispatching work request")
					workerChan <- work
				}()
			}
		}
	}()
}
