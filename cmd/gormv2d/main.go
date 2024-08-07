package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	dns := "root:12345#lxikm@tcp(127.0.0.1:3306)/dbv"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&LoginRecord{})
	if err != nil {
		panic(err)
	}
	t := time.Date(2024, time.June, 12, 9, 0, 0, 0, time.UTC)
	info := &LoginRecord{
		UID:   130,
		Year:  t.Year(),
		Month: int(t.Month()),
	}
	UpdateLogin(info, t.Day())
	// updataSql := fmt.Sprintf("UPDATE user_login_record t SET t.days = t.days | (1 << %d) WHERE t.uid = %d   and t.month = %d   and t.year = %d;", t.Day(), info.UID, t.Month(), t.Year())
	// result := db.Exec(updataSql)
	// if result.Error != nil {
	// 	panic(result)
	// }
	// if result.RowsAffected != 0 {
	// 	return
	// }

	result := db.Model(info).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}, {Name: "year"}, {Name: "month"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"days": gorm.Expr(fmt.Sprintf("days | (1 << %d)", t.Day()))}),
	}).Create(info)

	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println(result, info.ID)
}

type LoginRecord struct {
	ID         uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UID        int32  `gorm:"column:uid" json:"uid"`
	Year       int    `gorm:"column:year" json:"year"`
	Month      int    `gorm:"column:month" json:"month"`
	Days       uint32 `gorm:"column:days" json:"days"`
	UpdateTime int64  `gorm:"column:update_time" json:"update_time"`
}

func (l LoginRecord) TableName() string {
	return "user_login_record"
}

func (b *LoginRecord) CountConsecutiveDays(startDay int) (int, bool) {
	count := 0
	if startDay == 0 {
		totalDay := DaysInCurrentMonth(b.Year, b.Month)
		startDay = totalDay
	}

	for day := startDay; day >= 0 && CheckLoginDay(b, day); day-- {
		count++
	}
	return count, startDay == count
}

func DaysInCurrentMonth(year int, m int) int {

	month := time.Month(m)
	// 计算下个月的第一天
	if m == 12 {
		year++
		//month = time.January
	}

	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	// 计算本月的最后一天，即下个月的第一天的前一天
	lastDayThisMonth := firstDayNextMonth.Add(-24 * time.Hour)
	// 返回本月的天数
	return lastDayThisMonth.Day()
}

// UpdateLogin 登陆记录
func UpdateLogin(user *LoginRecord, dayOfMonth int) {
	if dayOfMonth < 0 || dayOfMonth >= 31 {
		return
	}
	user.Days |= 1 << dayOfMonth
}

// CheckLoginDay 检测是否登陆
func CheckLoginDay(user *LoginRecord, dayOfMonth int) bool {
	if dayOfMonth < 0 || dayOfMonth >= 31 {
		return false
	}
	return user.Days&(1<<dayOfMonth) != 0
}

type LoginConsecutiveDayes struct {
	HasLogin         bool `json:"has_login"`
	ConsecutiveDayes int  `json:"consecutive_dayes"`
}
