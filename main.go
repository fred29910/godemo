package main

import (
	"fmt"
	"time"
)

func main() {
	// 获取当前时间
	now := time.Now()

	// 获取当前是当月的第几天
	dayOfMonth := now.Day()

	// 输出结果
	fmt.Printf("今天是本月的第 %d 天\n", dayOfMonth)
}
