package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	i, err := GetUserWithdrawSc(db, context.Background(), 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(i)
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
func GetUserWithdrawSc(db *gorm.DB, ctx context.Context, uid int) (int, error) {
	// 查询单个数据
	var withdraw_sc int
	query := "SELECT withdraw_sc FROM players WHERE uid = ?"

	// 使用 QueryRow 来查询单个数据
	err := db.DB().QueryRow(query, uid).Scan(&withdraw_sc)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return withdraw_sc, err
}
