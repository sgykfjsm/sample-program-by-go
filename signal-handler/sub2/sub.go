package sub2

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Sub2(done chan<- bool) {
	sig := make(chan os.Signal, 1)
	subDone := make(chan bool, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sig
		fmt.Printf("\nCatch the signal at Sub2: %v\n", s)
		subDone <- true
	}()

	fmt.Println("Waiting for signal at Sub2 ...")
	<-subDone
	fmt.Println("Exit from Sub2")

	done <- true
}
