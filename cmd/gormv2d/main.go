package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. 创建一个 int 类型的 channel
	ch := make(chan int, 3)

	// 2. 启动一个新的 goroutine
	go func() {
		// 向 channel 发送一些数据
		for i := 1; i <= 3; i++ {
			fmt.Printf("发送数据: %d\n", i)
			ch <- i
			time.Sleep(time.Second)
		}
		// 关闭 channel
		close(ch)
		fmt.Println("Channel 已关闭")
	}()

	// 3. 在主 goroutine 中，使用 for...range 循环来读取 channel 中的数据
	fmt.Println("开始从 channel 接收数据...")
	for v := range ch {
		fmt.Printf("接收到数据: %v\n", v)
	}

	// 4. 验证行为
	fmt.Println("Range 循环结束，程序退出。")
}