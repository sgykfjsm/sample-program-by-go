package sub1

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Sub1(done chan<- bool) {
	sig := make(chan os.Signal, 1)
	subDone := make(chan bool, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sig
		fmt.Printf("\nCatch the signal at Sub1: %v\n", s)
		subDone <- true
	}()

	fmt.Println("Waiting for signal at Sub1 ...")
	<-subDone
	fmt.Println("Exit from Sub1")

	done <- true
}
