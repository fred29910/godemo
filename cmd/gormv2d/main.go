package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ids := []int{1, 2, 3, 4}
	sqtest := fmt.Sprintf("sqlxasd in (%#v)", ids)
	fmt.Println(sqtest)
}
