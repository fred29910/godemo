package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	dns := "root:12345#lxikm@tcp(127.0.0.1:3306)/dbv"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
	// createDatas(db)
	data, err := calConsecutiveDayes(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}

func createDatas(db *gorm.DB) {
	for i := 0; i < 6; i++ {
		t := time.Date(2024, time.August, i+1, 9, 0, 0, 0, time.UTC)
		createData(db, t)
	}
}

// createData creates a new LoginRecord in the database with the given UID, year,
// and month, and updates the days field to include the current day. If a record
// with the same UID, year, and month already exists, the days field is updated
// to include the current day.
//
// Parameters:
// - db: a pointer to a gorm.DB object representing the database connection.
//
// Return type: None.
func createData(db *gorm.DB, t time.Time) {
	// t := time.Date(2024, time.January, 3, 9, 0, 0, 0, time.UTC)
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
	fmt.Println(result.RowsAffected, info.ID)
}

// calConsecutiveDayes calculates the consecutive days for a given database connection.
//
// Parameters:
// - db: a pointer to a gorm.DB object representing the database connection.
//
// Return type: None.
func calConsecutiveDayes(db *gorm.DB) (*LoginConsecutiveDayes, error) {
	uid := 130
	t := time.Now()
	var logs []LoginRecord
	if err := db.Model(&LoginRecord{}).Where("uid = ?", uid).Order("id desc").Find(&logs).Error; err != nil {
		return nil, err
	}
	data := &LoginConsecutiveDayes{
		HasLogin: len(logs) > 0,
	}

	if len(logs) == 0 {
		data.ConsecutiveDayes = 0
		return data, nil
	}

	countUserLoginConsecutive := func() int {
		month := t.Month()
		year := t.Year()
		dayOfMonth := t.Day()
		count := 0
		for i, logdata := range logs {
			if i == 0 {
				if logdata.Year < year {
					return 0
				}
				if logdata.Month < int(month) {
					return 0
				}
			}
			dom := dayOfMonth
			if i != 0 {
				dom = DaysInCurrentMonth(logdata.Year, logdata.Month)
			}

			countMonth, ok := logdata.CountConsecutiveDays(dom)
			count += countMonth
			if !ok {
				break
			}
		}
		return count
	}
	count := countUserLoginConsecutive()

	data.ConsecutiveDayes = count

	return data, nil
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
	if dayOfMonth < 1 || dayOfMonth > 31 {
		return false
	}
	return user.Days&(1<<dayOfMonth) != 0
}

type LoginConsecutiveDayes struct {
	HasLogin         bool `json:"has_login"`
	ConsecutiveDayes int  `json:"consecutive_dayes"`
}
