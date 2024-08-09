package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	dsn := "root:12345#lxikm@tcp(localhost:3307)/dbv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name         string
	Age          sql.NullInt64
	Birthday     *time.Time
	Email        string  `gorm:"type:varchar(100);unique_index"`
	Role         string  `gorm:"size:255"`        // set field size to 255
	MemberNumber *string `gorm:"unique;not null"` // set member number to unique and not null
	Num          int     `gorm:"AUTO_INCREMENT"`  // set num to auto incrementable
	Address      string  `gorm:"index:addr"`      // create index with name `addr` for address
	IgnoreMe     int     `gorm:"-"`               // ignore this field
}
