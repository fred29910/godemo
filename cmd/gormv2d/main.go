package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSslMode := os.Getenv("DB_SSLMODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSslMode, dbTimezone)

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

	// "deuv:2134#lxikmv@tcp(127.0.0.1:3306)/dbv?parseTime=true&loc=Local"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
	err = db.Migrator().AutoMigrate(&MttAdminFee{})

	if err != nil {
		panic(err)
	}

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
	err = db.Create(&data).Error

	if err != nil {
		panic(err)
	}
}

type MttAdminFee struct {
	ID               uint64    `gorm:"column:id;type:bigserial;primaryKey;autoIncrement:true;comment:The squence id" json:"id"`                                              // The squence id
	PlayerID         uint32    `gorm:"column:player_id;type:integer;not null;uniqueIndex:idx_player_id_mtt_game_id,priority:1;comment:The Player ID" json:"player_id"`       // The Player ID
	MttGameID        uint32    `gorm:"column:mtt_game_id;type:integer;not null;uniqueIndex:idx_player_id_mtt_game_id,priority:2;comment:The mtt game id" json:"mtt_game_id"` // The mtt game id
	AdminFee         *float64  `gorm:"column:admin_fee;type:decimal;not null;default:0.00;comment:The admin fee" json:"admin_fee"`                                           // The admin fee
	MttGameStartTime time.Time `gorm:"column:mtt_game_start_time;type:timestamp;not null;comment:The mtt game start time" json:"mtt_game_start_time"`                        // The mtt game start time
	MttGameEndTime   time.Time `gorm:"column:mtt_game_end_time;type:timestamp;not null;comment:The mtt game end time" json:"mtt_game_end_time"`                              // The mtt game end time
	CreateTime       time.Time `gorm:"column:create_time;type:timestamp;not null;comment:The event create time" json:"create_time"`                                          // The event create time
	UpdateTime       time.Time `gorm:"column:update_time;type:timestamp;not null;comment:The event update time" json:"update_time"`                                          // The event update time
}
