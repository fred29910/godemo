package main

import (
	"fmt"
	"time"
)

func main() {
	// 获取当前时间
	now := time.Now()

	usWestTime(now)
	pvTime := time.Date(2023, time.December, 3, 3, 0, 9, 9, time.UTC)
	usWestTime(pvTime)
}

func usWestTime(now time.Time) {
	// 获取太平洋时区的时间
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// 将当前时间转化为太平洋时区时间
	pacificTime := now.In(loc)

	zoneName, offset := pacificTime.Zone()
	// 获取并打印当前时间及其时区信息
	fmt.Println("Current Time in Pacific Time Zone:", pacificTime)
	fmt.Println("Time Zone Name:", zoneName)
	fmt.Println("UTC Offset in Seconds:", offset)
}
