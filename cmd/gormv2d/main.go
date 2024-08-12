package main

import (
	"fmt"
	"time"
)

func main() {
	// 获取当前时间
	summerDay := time.Date(2024, 8, 11, 6, 0, 0, 9, time.UTC)

	usWestTime(summerDay)
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
	fmt.Println("\nCurrent Time in Pacific Time Zone:", pacificTime)
	fmt.Println("Time Zone Name:", zoneName)
	fmt.Println("UTC Offset in Seconds:", offset)

	fmt.Println("\nmonth:", pacificTime.Month())
	fmt.Println("day:", pacificTime.Day())
	fmt.Println("hour:", pacificTime.Hour())
}
