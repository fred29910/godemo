package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	s := make([]int, 2, 4) // len=2, cap=4
	s[0], s[1] = 1, 2
	fmt.Printf("Before append: len=%d, cap=%d, addr=%p\n", len(s), cap(s), s)

	c := append(s, 3) // 不超过 cap=4, 直接使用原数组
	fmt.Printf("After  append: len=%d, cap=%d, addr=%p %v\n", len(c), cap(c), c, c)
	fmt.Printf("After  append: len=%d, cap=%d, addr=%p %v\n", len(s), cap(s), s, s)

	d := append(s, 4, 5, 8, 5) // 超出 cap=4, 需要分配新数组
	fmt.Printf("After 2nd append: len=%d, cap=%d, addr=%p %v\n", len(d), cap(d), d, d)
	fmt.Printf("After 2nd append: len=%d, cap=%d, addr=%p %v\n", len(s), cap(s), s, s)

}
