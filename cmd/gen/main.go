package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func main() {

	// genGormDb("/home/fred/workspace/wptglobal/clubwpt-backend/build/database/db_coin_test/", "internal/coin/query")
	// genGormDb("/home/fred/workspace/wptglobal/clubwpt-backend/build/database/db_service_test/", "internal/service/query")
	genGormDb("/home/fred/workspace/wptglobal/clubwpt-backend/build/database/dzpk_static_test/", "internal/dzpk_static_test/query")
}

func genGormDb(sqlPath string, outPath string) {

	// dns := "root:12345#lxikm@tcp(127.0.0.1:3306)/dbv?parseTime=true&loc=Local"
	// gormdb, err := gorm.Open(mysql.Open(dns), &gorm.Config{
	// 	Logger: newLogger,
	// })
	// sqlPath := "/home/fred/workspace/wptglobal/clubwpt-backend/build/database/db_coin_test/"
	var sqlfiles []string
	// 遍历目录
	err := filepath.Walk(sqlPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".sql") {
			return nil
		}
		sqlfiles = append(sqlfiles, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	gormdb, err := gorm.Open(rawsql.New(rawsql.Config{
		//SQL:      rawsql,                      //create table sql
		FilePath: sqlfiles,
	}))

	if err != nil {
		log.Fatal(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           outPath,                                                            //"internal/query"
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
