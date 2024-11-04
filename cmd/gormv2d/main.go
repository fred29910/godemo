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
			// ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful: false, // Disable color
		},
	)

	dns := "root:12345#lxikm@tcp(127.0.0.1:3307)/dbv?parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	// db.Model(&PromoCode{}).Create(&PromoCode{
	// 	Status: 1,
	// 	Code:   "test",
	// 	Info:   "test",
	// })

	var promoCodes []*PromoCode

	query := db.Model(&PromoCode{})

	query.Where("status = ?", 1)
	query.Where("info like ?", `%test%`)

	query.Where("created_at < ?", time.Now())
	err = query.Find(&promoCodes).Error

	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("%+v", promoCodes))
}

type PromoCode struct {
	gorm.Model
	Status int    `json:"status" gorm:"column:status"`
	Code   string `json:"code" gorm:"column:code"`
	Info   string `json:"info" gorm:"column:info"`
}

const PromoCodeTableName = "promo_code"

func (PromoCode) TableName() string {
	return PromoCodeTableName
}

// type PromoCodeDao struct {
// 	Status    int        `json:"status" gorm:"column:status"`
// 	Code      string     `json:"code" gorm:"column:code"`
// 	Info      string     `json:"info" gorm:"column:info"`
// 	ID        uint       `json:"id" gorm:"column:id"`
// 	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
// }
