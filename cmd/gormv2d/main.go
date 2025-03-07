package main

import (
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

	dns := "root:12345#lxikm@tcp(127.0.0.1:3306)/dbv?parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
	// err = db.Migrator().AutoMigrate(&ScFinanceReport{})

	// if err != nil {
	// 	panic(err)
	// }

	// attr := map[string]interface{}{
	// 	"test":  "test",
	// 	"test1": "test1",
	// 	"test2": "test2",
	// }

	// atrJson, _ := json.Marshal(attr)
	// db.Model(&PromoCode{}).Create(&PromoCode{
	// 	Status: 1,
	// 	Code:   "test",
	// 	Info:   "test",

	// 	Attributes: datatypes.JSON(atrJson),
	// })
	adminf := 1.01

	data := MttAdminFee{
		ID:               1,
		PlayerID:         1,
		MttGameID:        1,
		AdminFee:         &adminf,
		MttGameStartTime: time.Now(),
		MttGameEndTime:   time.Now(),
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	err = db.Table("mtt_admin_fees_202503").Create(&data).Error

	if err != nil {
		panic(err)
	}
}

type MttAdminFee struct {
	ID               uint64    `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true;comment:The squence id" json:"id"`                                             // The squence id
	PlayerID         uint32    `gorm:"column:player_id;type:int(11) unsigned;not null;uniqueIndex:idx_player_id_mtt_game_id,priority:1;comment:The Player ID" json:"player_id"`       // The Player ID
	MttGameID        uint32    `gorm:"column:mtt_game_id;type:int(11) unsigned;not null;uniqueIndex:idx_player_id_mtt_game_id,priority:2;comment:The mtt game id" json:"mtt_game_id"` // The mtt game id
	AdminFee         *float64  `gorm:"column:admin_fee;type:decimal(18,2);not null;default:0.00;comment:The admin fee" json:"admin_fee"`                                              // The admin fee
	MttGameStartTime time.Time `gorm:"column:mtt_game_start_time;type:datetime;not null;comment:The mtt game start time" json:"mtt_game_start_time"`                                  // The mtt game start time
	MttGameEndTime   time.Time `gorm:"column:mtt_game_end_time;type:datetime;not null;comment:The mtt game end time" json:"mtt_game_end_time"`                                        // The mtt game end time
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:The event create time" json:"create_time"`                                                    // The event create time
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:The event update time" json:"update_time"`                                                    // The event update time
}
