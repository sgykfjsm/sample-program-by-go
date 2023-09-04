package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 1. キャンセル伝搬
	ctx, cancel := context.WithCancel(context.Background())
	go operation1(ctx)

	// タイムアウト
	timeoutCtx, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	go operation2(timeoutCtx)
	cancel2()

	deadlineCtx, cancel3 := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	go operation2(deadlineCtx)
	time.Sleep(3 * time.Second)
	cancel3()

	// 値の保存と取得
	valueCtx := context.WithValue(ctx, "userID", 1)
	go operation3(valueCtx)
	valueCtx2 := context.WithValue(ctx, "userID", 2)
	go operation3(valueCtx2)

	// キャンセル実行
	// すべての処理がキャンセルされる。
	cancel()

	time.Sleep(5 * time.Second)
	fmt.Println("good bye")
}

func operation1(ctx context.Context) {
	// 非同期処理の制御、エラーハンドリング
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("operation1 completed")
	case <-ctx.Done():
		fmt.Println("operation1 cancelled:", ctx.Err())
	}
}

func operation2(ctx context.Context) {
	// APIの一貫性、トレーシング
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("operation2 completed")
	case <-ctx.Done():
		fmt.Println("operation2 cancelled:", ctx.Err())
	}
}

func operation3(ctx context.Context) {
	// 値の取得
	userID := ctx.Value("userID").(int)
	fmt.Println("operation3 started for userID: ", userID)
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("operation3 completed")
	case <-ctx.Done():
		fmt.Println("operation3 cancelled:", ctx.Err())
	}
}
