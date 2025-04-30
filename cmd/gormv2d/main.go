package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// 获取当前时间
	// summerDay := time.Date(2024, 8, 11, 6, 0, 0, 9, time.UTC)

	// usWestTime(summerDay)
	loc, _ := time.LoadLocation("America/Los_Angeles")
	JaStartTimestamp := time.Date(2024, time.January, 1, 0, 0, 0, 0, loc)
	usWestTime(JaStartTimestamp)

	freStartTime := time.Date(2024, time.February, 1, 0, 0, 0, 0, loc)
	usWestTime(freStartTime)

	freEndTime := time.Date(2024, time.March, 1, 0, 0, 0, 0, loc)
	usWestTime(freEndTime)

	logrus.Infof("%d - %d", JaStartTimestamp.Unix(), freStartTime.Unix())
	logrus.Infof("%d - %d", freStartTime.Unix(), freEndTime.Unix())

	// pvTime := time.Unix(1738464000, 0)
	// usWestTime(pvTime)
	// usWestTime(time.Now())
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
