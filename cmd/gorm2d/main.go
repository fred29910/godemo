package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:12345#lxikm@tcp(127.0.0.1:3306)/coin?charset=utf8mb4&parseTime=True&loc=Local"
	DbEngin, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// DbEngin.Migrator().AutoMigrate(&MiniGameBetPreferenceTag{})
	// enable := int32(1)
	// DbEngin.Create(&MiniGameBetPreferenceTag{
	// 	GameID:           1,
	// 	Status:           &enable,
	// 	TagNameCn:        "xasdasda",
	// 	TagNameEn:        "xasda",
	// 	MatchCondition:   "xasdfasdda",
	// 	Option:           "xasdasa",
	// 	WeightPerHit:     21,
	// 	TriggerThreshold: 90,
	// })

	var dsin []MiniGameBetPreferenceTag
	if err := DbEngin.Find(&dsin).Error; err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(dsin)
	fmt.Printf("%v\n", string(bs))
}

type MiniGameBetPreferenceTag struct {
	ID               uint64         `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true" json:"id"`
	GameID           uint32         `gorm:"column:game_id;type:int(11) unsigned;not null;index:idx_game_id,priority:1;comment:所属游戏" json:"game_id"`    // 所属游戏
	Status           *int32         `gorm:"column:status;type:tinyint(4);not null;default:2;comment:enable:1, disable: 2" json:"status"`               // enable:1, disable: 2
	TagNameCn        string         `gorm:"column:tag_name_cn;type:varchar(255);not null;comment:标签名称（中文），用于前端展示" json:"tag_name_cn"`                  // 标签名称（中文），用于前端展示
	TagNameEn        string         `gorm:"column:tag_name_en;type:varchar(255);not null;comment:标签名称（英文），用于前端展示" json:"tag_name_en"`                  // 标签名称（英文），用于前端展示
	MatchCondition   string         `gorm:"column:match_condition;type:varchar(255);not null;comment:标签所对应的下注行为匹配条件，不支持后台修改" json:"match_condition"`   // 标签所对应的下注行为匹配条件，不支持后台修改
	Option           string         `gorm:"column:option;type:varchar(255);not null;comment:偏好值,数组，目前只支持or" json:"option"`                             // 偏好值,数组，目前只支持or
	WeightPerHit     int32          `gorm:"column:weight_per_hit;type:int(11);not null;comment:玩家每次命中该标签行为时累加的分数" json:"weight_per_hit"`               // 玩家每次命中该标签行为时累加的分数
	TriggerThreshold int32          `gorm:"column:trigger_threshold;type:int(11);not null;comment:玩家最近100局中累计得分达到该数值才进入候选范围" json:"trigger_threshold"` // 玩家最近100局中累计得分达到该数值才进入候选范围
	Notes            *string        `gorm:"column:notes;type:text;comment:备注信息" json:"notes"`                                                          // 备注信息
	CreatedAt        *UTCTime       `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:记录创建时间" json:"created_at"`       // 记录创建时间
	UpdatedAt        *UTCTime       `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:记录最后更新时间" json:"updated_at"`     // 记录最后更新时间
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"delete_at,omitempty"`
}

func (m *MiniGameBetPreferenceTag) TableName() string {
	return "mini_game_bet_preference_tag"
}

type UTCTime struct {
	time.Time
}

func (t UTCTime) Value() (driver.Value, error) {
	// 如果时间是零值，则存入 NULL
	if t.Time.IsZero() {
		return nil, nil
	}
	// 确保存入的是UTC时间
	return t.Time.UTC(), nil
}

// Scan 实现 sql.Scanner 接口，在从数据库读取时被调用
func (t *UTCTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{} // or handle as you wish
		return nil
	}

	var err error
	switch v := value.(type) {
	case time.Time:
		// 数据库驱动已将其解析为 time.Time
		t.Time = v.UTC()
	case []byte:
		// 可能是字符串或时间格式的字节数组
		// time.RFC3339Nano 是一个常见的数据库时间格式
		t.Time, err = time.Parse(time.RFC3339Nano, string(v))
		if err != nil {
			// 尝试其他可能的格式
			t.Time, err = time.Parse("2006-01-02 15:04:05", string(v))
		}
		t.Time = t.Time.UTC()
	case string:
		t.Time, err = time.Parse(time.RFC3339Nano, v)
		if err != nil {
			t.Time, err = time.Parse("2006-01-02 15:04:05", v)
		}
		t.Time = t.Time.UTC()
	default:
		return fmt.Errorf("cannot scan type %T into UTCTime: %v", value, value)
	}

	return err
}
