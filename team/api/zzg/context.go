package main

import (
	"context"
	"fmt"
	"time"
)

// 模拟一个最小执行时间的阻塞函数
func inc(a int) int {
	res := a + 1
	time.Sleep(1 * time.Second)
	return res
}

// 向外部提供的阻塞接口
func Add1(ctx context.Context, a int) int {
	res := 0
	for i := 0; i < a; i++ {
		res = inc(res)
		select {
		case <-ctx.Done():
			return -1
		default:
		}
	}

	return res
}

func Add2(ctx context.Context, b int) int {
	res := 0
	for i := 0; i < b; i++ {
		res = inc(res)
		select {
		case <-ctx.Done():
			return -1
		default:
		}
	}
	return res
}
func main() {
	{
		// 使用开放的 API 计算 a+b
		timeout := 2 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		res := Add1(ctx, 1)
		fmt.Printf("result: %d\n", res)
	}
	{
		// 使用开放的 API 计算 a+b
		timeout := 2 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		res := Add1(ctx, 3)
		fmt.Printf("result: %d\n", res)
	}
	{
		// 手动取消
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		res := Add2(ctx, 1)
		fmt.Printf("result: %d\n", res)
	}
}
