package main

import (
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dxs := struct {
		Isat time.Time `gorm:"column:isat" json:"isat"`
	}{
		Isat: time.Now(),
	}
	bs, _ := json.Marshal(dxs)

	println(string(bs))

}
