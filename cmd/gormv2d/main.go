package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if daysInCurrentMonth(2023, 12) != 31 {
		panic("get error")
	}
	if daysInCurrentMonth(2023, 4) != 30 {
		panic("get error")
	}
}

func daysInCurrentMonth(year int, m uint8) int {

	month := time.Month(m)
	// 计算下个月的第一天
	if m == 12 {
		year++
		//month = time.January
	}

	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	// 计算本月的最后一天，即下个月的第一天的前一天
	lastDayThisMonth := firstDayNextMonth.Add(-24 * time.Hour)
	// 返回本月的天数
	return lastDayThisMonth.Day()
}
