package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func sleepFunc() error {
	time.Sleep(3 * time.Second)
	return nil
}

func ContextWrapper(ctx context.Context, f func() error) error {
	fmt.Println("Please Ctrl+C after a while")

	errChan := make(chan error, 1)
	go func() {
		errChan <- f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func main() {
	// trap Ctrl+C and call cancel on the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case s := <-c:
			fmt.Printf("\nCatch the Signal %q\n", s.String())
			cancel()
		}
	}()

	if err := ContextWrapper(ctx, sleepFunc); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Task is finished without error")
}
