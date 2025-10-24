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
	concurrentJobs := 1000 // 模拟 1000 个并发作业

	wg.Add(concurrentJobs)

	fmt.Println("开始并发测试，启动 1000 个 Goroutines...")

	// 2. 模拟并发读写
	for i := 0; i < concurrentJobs; i++ {
		go func(jobID int) {
			defer wg.Done()

			// 我们混合读写操作

			if jobID%10 == 0 {
				// 大约 10% 的作业只读取当前状态
				// f.Current() 会获取 stateMu.RLock()
				state := f.Current()
				_ = state // 消除 "unused" 警告
			} else if jobID%3 == 0 {
				// 尝试 "run"
				// f.Event() 会获取 stateMu.Lock()
				// FSM 内部的锁会处理冲突。
				// 如果状态不是 "pending"，f.Event 会安全地返回错误。
				_ = f.Event(context.Background(), "run")
			} else if jobID%3 == 1 {
				// 尝试 "stop"
				_ = f.Event(context.Background(), "stop")
			} else {
				// 尝试 "reset"
				_ = f.Event(context.Background(), "reset")
			}
		}(i)
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
