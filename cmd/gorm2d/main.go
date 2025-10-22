package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

func main() {

	var datas []UserHandInfo
	for i := 0; i < 100; i++ {
		datas = append(datas, UserHandInfo{
			HandNum: 1000,
			Cnt:     time.Now().Unix(),
			Options: "1,2,3,4,5,6,7,302",
		})
	}

	bs, _ := json.Marshal(datas)
	fmt.Printf("data len %v \n", humanize.Bytes(uint64(len(bs))))
}

type UserHandInfo struct {
	HandNum int64  `gorm:"column:hand_num" json:"hand_num"`
	Cnt     int64  `gorm:"column:cnt" json:"record_time"`
	Options string `gorm:"column:options" json:"options"`
}
