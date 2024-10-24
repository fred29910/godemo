package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 设置日志级别为详细模式
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	dns := "root:12345#lxikm@tcp(127.0.0.1:3307)/dbv"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
	var data string
	err = db.Table(TableName).Select("ip").Where("id = ?", 1).Pluck("ip", &data).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}

const TableName = "coin_white_ip"
