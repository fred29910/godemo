package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbData := MysqlConf{
		Ip:           "127.0.0.1",
		Pwd:          "12345#lxikm",
		User:         "root",
		Port:         3306,
		Database:     "texas",
		LoadInterval: 5,
	}
	dataDns := dbData.User + ":" + dbData.Pwd + "@tcp(" + dbData.Ip + ":" + strconv.Itoa(dbData.Port) + ")/" + dbData.Database + "?charset=utf8&loc=Local&readTimeout=15s&writeTimeout=15s"

	driverName := "mysql"
	dbs, err := sql.Open(driverName, dataDns)
	if err != nil {
		panic(err)
	}
	defer dbs.Close()

	result, err := dbs.Exec("insert into players(uid, withdraw_sc) value (43, 1) on duplicate key update withdraw_sc = withdraw_sc+1;")
	if err != nil {
		log.Fatal("cat not get ", err.Error())
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Fatal("cat not get count", err.Error())
		return
	}
	log.Printf("result %v", count)
}

type MysqlConf struct {
	Ip           string `xml:",attr"`
	Port         int    `xml:",attr"`
	User         string `xml:",attr"`
	Pwd          string `xml:",attr"`
	Database     string `xml:",attr"`
	LoadInterval int    `xml:",attr"`
}
