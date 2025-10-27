package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/looplab/fsm" // 确保已安装: go get github.com/looplab/fsm
)

func main() {
	// 1. 定义一个状态机实例
	f := fsm.NewFSM(
		"pending", // 初始状态
		fsm.Events{
			// 定义循环事件，以便并发测试可以持续触发
			{Name: "run", Src: []string{"pending"}, Dst: "running"},
			{Name: "stop", Src: []string{"running"}, Dst: "stopped"},
			{Name: "reset", Src: []string{"stopped"}, Dst: "pending"},
		},
		fsm.Callbacks{},
	)

	var wg sync.WaitGroup

	// 2. 模拟并发读写
	for i := 0; i < 100; i++ {
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(jobID int) {
				defer wg.Done()

				// 我们混合读写操作

				// 如果状态不是 "pending"，f.Event 会安全地返回错误。
				_ = f.Event(context.Background(), "run")

			}(i)
		}
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(jobID int) {
				defer wg.Done()

				// 我们混合读写操作

				// 尝试 "stop"
				_ = f.Event(context.Background(), "stop")

			}(i)
		}
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(jobID int) {
				defer wg.Done()

				// 我们混合读写操作

				// 尝试 "reset"
				_ = f.Event(context.Background(), "reset")

			}(i)
		}
	}

	// 3. 等待所有 goroutine 完成
	wg.Wait()

	fmt.Println("并发测试完成。")
	fmt.Println("最终状态:", f.Current()) // 最终状态是随机的，取决于最后一个写入操作
	fmt.Println("\n--- 如何验证 ---")
	fmt.Println("1. 确保已安装 'github.com/looplab/fsm'")
	fmt.Println("   (go get github.com/looplab/fsm)")
	fmt.Println("2. 将此代码保存为 main.go")
	fmt.Println("3. 运行 'go run -race main.go'")
	fmt.Println("4. 如果程序正常结束且未报告 'DATA RACE'，则证明 FSM 是并发安全的。")
}
