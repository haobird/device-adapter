package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 测试
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int, 10)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("gen---done--", n)
					return // returning not to leak the goroutine
				case dst <- n:
					fmt.Println("gen---dst--", n)
					n++

				}
			}
		}()
		fmt.Println("gen---return--", n)
		return dst
	}

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}

	gen(ctx)

	time.AfterFunc((1 * time.Minute), func() {
		fmt.Println("执行after")
		cancel()
	})
	keepAlive()
}

func keepAlive() {
	//合建chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("总进程服务启动完成")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)
}
