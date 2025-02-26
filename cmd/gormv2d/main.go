package main

import (
	"log"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func main() {

	// dns := "root:12345#lxikm@tcp(127.0.0.1:3306)/dbv?parseTime=true&loc=Local"
	// gormdb, err := gorm.Open(mysql.Open(dns), &gorm.Config{
	// 	Logger: newLogger,
	// })
	gormdb, err := gorm.Open(rawsql.New(rawsql.Config{
		//SQL:      rawsql,                      //create table sql
		FilePath: []string{
			//"./sql/user.sql", // create table sql file
			"/home/fred/workspace/wptglobal/clubwpt-backend/build/database/db_coin_test/rewards_cfg.sql", // create table sql file directory
		},
	}))

	if err != nil {
		log.Fatal(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "internal/query",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		FieldNullable:     true,                                                               // generate pointer at struct
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		WithUnitTest:      true,
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions

	// g.ApplyBasic(
	// 	// Generate struct `User` based on table `users`
	// 	g.GenerateModel("users"),

	// 	// Generate struct `Employee` based on table `users`
	// 	g.GenerateModelAs("users", "Employee"),

	// 	// Generate struct `User` based on table `users` and generating options
	// 	g.GenerateModel("users", gen.FieldIgnore("address"), gen.FieldType("id", "int64")),

	// 	// Generate struct `Customer` based on table `customer` and generating options
	// 	// customer table may have a tags column, it can be JSON type, gorm/gen tool can generate for your JSON data type
	// 	g.GenerateModel("customer", gen.FieldType("tags", "datatypes.JSON")),
	// )
	g.ApplyBasic(
		// Generate structs from all tables of current database
		g.GenerateAllTable()...,
	)
	// Generate the code
	g.Execute()
}
