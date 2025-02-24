package main

import (
	"log"
	"os"
	"time"

	"gorm.io/datatypes"
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

	dataRp := &ScFinanceReport{
		Hour:                time.Now(),
		StartTime:           time.Now(),
		EndTime:             time.Now(),
		ScTotalEnding:       100,
		ScRedeemableEnding:  100,
		UsdPurchased:        100,
		ScWasRedeemed:       100,
		ScRedeemableCreated: 100,

		ScRedeemableRemoved: map[uint32]uint64{1: 232, 2: 232},
		ScTotalCreated:      map[uint32]uint64{1: 232, 2: 232},
		ScTotalRemoved:      nil,
	}

	err = db.Model(&ScFinanceReport{}).Create(&dataRp).Error
	if err != nil {
		panic(err)
	}

}

type PromoCode struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Status int    `json:"status" gorm:"column:status"`
	Code   string `json:"code" gorm:"column:code"`
	Info   string `json:"info" gorm:"column:info"`

	Attributes datatypes.JSON `json:"attributes" gorm:"column:attributes"`
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

type ScFinanceReport struct {
	ID        uint64    `json:"id"`
	Hour      time.Time `gorm:"UNIQUE_INDEX:once_id_hour" json:"hour"`
	CreatedAt time.Time `json:"created_at"` // 資料紀錄時間
	UpdatedAt time.Time `json:"updated_at"` // 資料修改時間

	StartTime           time.Time         `json:"start_time"` // 开始时间
	EndTime             time.Time         `json:"end_time"`   // 结束时间
	ScTotalEnding       uint64            `json:"sc_total_ending"`
	ScRedeemableEnding  uint64            `json:"sc_redeemable_ending"`
	UsdPurchased        uint64            `json:"usd_purchasesd"`
	ScWasRedeemed       uint64            `json:"sc_was_redeemed,omitempty"`
	ScRedeemableCreated uint64            `json:"sc_redeemable_created,omitempty"`
	ScRedeemableRemoved map[uint32]uint64 `json:"sc_redeemable_removed,omitempty" gorm:"serializer:json"`
	ScTotalCreated      map[uint32]uint64 `json:"sc_total_created,omitempty" gorm:"serializer:json"`
	ScTotalRemoved      map[uint32]uint64 `json:"sc_total_removed,omitempty" gorm:"serializer:json"`

	Rake      uint64 `json:"rake"`
	SystemEat uint64 `json:"system_eat"`
}

func (s *ScFinanceReport) TableName() string {
	return "sc_finance_report"
}
