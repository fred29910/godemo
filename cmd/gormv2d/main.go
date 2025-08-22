package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// BatchQueue 是一个支持泛型的并发安全批量处理队列
type BatchQueue[T any] struct {
	// 配置项
	maxSize    int           // 批处理大小阈值
	maxTimeout time.Duration // 批处理时间阈值

	// 内部状态
	inputChan  chan T             // 用于接收数据的通道
	batchFunc  func([]T)          // 处理批次数据的回调函数
	ctx        context.Context    // 控制生命周期的 context
	cancelFunc context.CancelFunc // 用于发出关闭信号
	wg         sync.WaitGroup     // 用于等待 processor goroutine 结束
}

// NewBatchQueue 创建一个新的批量处理队列实例
//
// 参数:
//
//	maxSize:    队列满多少笔数据后触发批量处理 (e.g., 100)
//	maxTimeout: 等待多久后触发批量处理 (e.g., 15 * time.Second)
//	batchFunc:  处理一个批次数据的函数
func NewBatchQueue[T any](maxSize int, maxTimeout time.Duration, batchFunc func([]T)) (*BatchQueue[T], error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("maxSize must be greater than 0")
	}
	if maxTimeout <= 0 {
		return nil, fmt.Errorf("maxTimeout must be a positive duration")
	}
	if batchFunc == nil {
		return nil, fmt.Errorf("batchFunc cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	q := &BatchQueue[T]{
		maxSize:    maxSize,
		maxTimeout: maxTimeout,
		batchFunc:  batchFunc,
		// 使用一个带缓冲的 channel，提高吞吐量
		// 缓冲大小可以根据实际情况调整，这里设为 maxSize
		inputChan:  make(chan T, maxSize),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	// 启动后台处理 goroutine
	q.wg.Add(1)
	go q.processor()

	return q, nil
}

// Add 向队列中添加一个数据项。这是并发安全的。
func (q *BatchQueue[T]) Add(item T) {
	// 使用 select 避免在队列关闭后写入导致 panic
	select {
	case q.inputChan <- item:
	case <-q.ctx.Done():
		log.Printf("Warning: Queue is closed. Item was not added.")
	}
}

// Shutdown 优雅地关闭队列
// 它会等待所有已在队列中的数据被处理完毕
func (q *BatchQueue[T]) Shutdown() {
	// 发出关闭信号
	q.cancelFunc()
	// 等待 processor goroutine 完全退出
	q.wg.Wait()
}

// processor 是后台运行的核心处理逻辑
func (q *BatchQueue[T]) processor() {
	defer q.wg.Done()
	// 关闭 inputChan，这样 range 循环才能在 Shutdown 后正常退出
	// 注意：这需要在 Shutdown() 被调用后，且没有新的 Add() 调用时才安全
	// 在我们的设计中，cancelFunc 会让 select 退出，所以这里可以安全地 close
	// 不过更安全的做法是在 shutdown 里 close
	defer close(q.inputChan)

	// 初始化一个缓冲区，容量预设为 maxSize 提高效率
	buffer := make([]T, 0, q.maxSize)
	// 初始化一个定时器
	timer := time.NewTimer(q.maxTimeout)
	// 确保 timer 在未使用时被正确停止，避免资源泄露
	defer timer.Stop()

	for {
		select {
		// 1. 接收到新的数据项
		case item, ok := <-q.inputChan:
			if !ok {
				// inputChan 被关闭，说明即将退出
				// 处理缓冲区中最后剩余的数据
				if len(buffer) > 0 {
					log.Println("Input channel closed. Processing final batch...")
					q.batchFunc(buffer)
				}
				return
			}

			buffer = append(buffer, item)
			// 检查是否达到数量阈值
			if len(buffer) >= q.maxSize {
				log.Printf("Batch triggered by size: %d items\n", len(buffer))
				q.batchFunc(buffer)
				// 清空缓冲区
				buffer = make([]T, 0, q.maxSize)
				// 重置定时器
				// 先 stop 确保 timer 的 channel 是空的，这是 Go timer 的标准实践
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(q.maxTimeout)
			}

		// 2. 定时器触发
		case <-timer.C:
			if len(buffer) > 0 {
				log.Printf("Batch triggered by timeout: %d items\n", len(buffer))
				q.batchFunc(buffer)
				// 清空缓冲区
				buffer = make([]T, 0, q.maxSize)
			}
			// 无论缓冲区是否为空，都需要重置定时器
			timer.Reset(q.maxTimeout)

		// 3. 接收到关闭信号
		case <-q.ctx.Done():
			// 处理缓冲区中最后剩余的数据
			if len(buffer) > 0 {
				log.Println("Shutdown signal received. Processing final batch...")
				q.batchFunc(buffer)
			}
			return
		}
	}
}

// main 函数：演示如何使用 BatchQueue
func main() {
	log.Println("Starting batch queue demonstration...")

	// 模拟的批量更新函数
	// 在真实场景中，这里会执行数据库批量插入或调用批量更新 API
	batchUpdateFunc := func(items []int) {
		fmt.Printf("-----> Executing batch update with %d items: %v\n", len(items), items)
		// 模拟网络或数据库延迟
		time.Sleep(100 * time.Millisecond)
	}

	// 创建队列：满 10 笔或 3 秒执行一次
	// 为了方便演示，这里使用较小的值
	queue, err := NewBatchQueue(10, 3*time.Second, batchUpdateFunc)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}
	// 确保在 main 函数退出时优雅关闭队列
	defer queue.Shutdown()

	// --- 场景1: 快速添加数据，触发数量阈值 ---
	log.Println("\n--- Scenario 1: Triggering by size limit ---")
	for i := 1; i <= 12; i++ {
		fmt.Printf("Adding item %d\n", i)
		queue.Add(i)
		time.Sleep(50 * time.Millisecond) // 模拟快速连续的请求
	}

	// 等待一下，让上面的批处理完成并打印日志
	time.Sleep(1 * time.Second)

	// --- 场景2: 慢速添加数据，触发时间阈值 ---
	log.Println("\n--- Scenario 2: Triggering by timeout ---")
	for i := 13; i <= 15; i++ {
		fmt.Printf("Adding item %d\n", i)
		queue.Add(i)
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("Waiting for timeout trigger...")
	// 等待超过3秒，让定时器触发
	time.Sleep(4 * time.Second)

	// --- 场景3: 关闭前处理剩余数据 ---
	log.Println("\n--- Scenario 3: Processing remaining items on shutdown ---")
	fmt.Println("Adding item 16")
	queue.Add(16)
	fmt.Println("Adding item 17")
	queue.Add(17)

	log.Println("Shutting down the queue...")
	// Shutdown 会在退出前处理掉 16 和 17
}
