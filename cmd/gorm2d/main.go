package main

import (
	"godemo/internal/coin/query"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))

	query.SetDefault(gormdb)

	qtx := query.Q.Begin()

	qtx.Agreement.Create()
	qtx.Commit()

}
