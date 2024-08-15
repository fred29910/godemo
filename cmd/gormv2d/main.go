package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	dsn := "root:12345#lxikm@tcp(localhost:3307)/coin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()
	GetUserRegisterRewardTime(db, context.Background(), 1)
}

func GetUserRegisterRewardTime(db *gorm.DB, ctx context.Context, uid int) (int64, error) {
	var logInRds = struct {
		RewardTime int64 `gorm:"column:register_reward_time" json:"register_reward_time"`
	}{}

	err := db.Table("dtb_user_main").Where("id = ?", uid).Select("register_reward_time").First(&logInRds).Error
	if err != nil {
		return 0, err
	}
	return logInRds.RewardTime, nil
}
