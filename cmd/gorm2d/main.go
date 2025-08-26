package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:12345#lxikm@tcp(127.0.0.1:3306)/coin_92_test?charset=utf8mb4&parseTime=True&loc=Local"
	DbEngin, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	var data []*MiniGameBetPreferenceTag
	if err := DbEngin.Table("mini_game_bet_preference_tag").Find(&data).Error; err != nil {
		log.Fatalf("cat not query data %v", err)
	}
	bsm, _ := json.Marshal(data)
	fmt.Printf("%v", string(bsm))

}

type MiniGameBetPreferenceTag struct {
	ID               uint64     `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true" json:"id"`
	GameID           uint32     `gorm:"column:game_id;type:int(11) unsigned;not null;index:idx_game_id,priority:1;comment:所属游戏" json:"game_id"`    // 所属游戏
	Status           *int32     `gorm:"column:status;type:tinyint(4);not null;default:2;comment:enable:1, disable: 2" json:"status"`               // enable:1, disable: 2
	TagNameCn        string     `gorm:"column:tag_name_cn;type:varchar(255);not null;comment:标签名称（中文），用于前端展示" json:"tag_name_cn"`                  // 标签名称（中文），用于前端展示
	TagNameEn        string     `gorm:"column:tag_name_en;type:varchar(255);not null;comment:标签名称（英文），用于前端展示" json:"tag_name_en"`                  // 标签名称（英文），用于前端展示
	MatchCondition   string     `gorm:"column:match_condition;type:varchar(255);not null;comment:标签所对应的下注行为匹配条件，不支持后台修改" json:"match_condition"`   // 标签所对应的下注行为匹配条件，不支持后台修改
	Option           string     `gorm:"column:option;type:varchar(255);not null;comment:偏好值,数组，目前只支持or" json:"option"`                             // 偏好值,数组，目前只支持or
	WeightPerHit     int32      `gorm:"column:weight_per_hit;type:int(11);not null;comment:玩家每次命中该标签行为时累加的分数" json:"weight_per_hit"`               // 玩家每次命中该标签行为时累加的分数
	TriggerThreshold int32      `gorm:"column:trigger_threshold;type:int(11);not null;comment:玩家最近100局中累计得分达到该数值才进入候选范围" json:"trigger_threshold"` // 玩家最近100局中累计得分达到该数值才进入候选范围
	Notes            *string    `gorm:"column:notes;type:text;comment:备注信息" json:"notes"`                                                          // 备注信息
	CreatedAt        *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:记录创建时间" json:"created_at"`       // 记录创建时间
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:记录最后更新时间" json:"updated_at"`     // 记录最后更新时间
}
