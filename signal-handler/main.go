package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sgykfjsm/sample-program-by-go/signal-handler/sub1"
	"github.com/sgykfjsm/sample-program-by-go/signal-handler/sub2"
)

func main() {
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	sub1Done := make(chan bool, 1)
	sub2Done := make(chan bool, 1)

	go sub1.Sub1(sub1Done)
	go sub2.Sub2(sub2Done)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sig
		fmt.Printf("\nCatch the signal at main: %v\n", s)
		done <- true
	}()

	fmt.Println("Waiting for signal at main ...")
	<-done
	fmt.Println("Exit from main")

	<-sub1Done
	<-sub2Done
}
