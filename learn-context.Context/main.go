package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type myKey string

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := myKey("duration")
	ctx = context.WithValue(ctx, key, 2*time.Second)

	done := make(chan string)

	go func(ctx context.Context) {
		if duration, ok := ctx.Value(key).(time.Duration); ok {
			fmt.Printf("Sleep %s\n", duration.String())
			time.Sleep(duration)
			done <- "Done!"
		} else {
			log.Fatal("something wrong")
		}
	}(ctx)

	select {
	case <-ctx.Done():
		fmt.Println("Timeout!")
	case result := <-done:
		fmt.Println(result)
	}
}
