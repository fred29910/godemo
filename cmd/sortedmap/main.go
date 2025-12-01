package main

import (
	"fmt"
	"godemo/pkg/sortedmap"
)

func main() {
	cmpInt := func(a, b int) int {
		return a - b
	}

	sm := sortedmap.NewSortedMapByValue[string, int](cmpInt)

	sm.Put("a", 30)
	sm.Put("b", 10)
	sm.Put("c", 20)
	sm.Put("d", 20)

	fmt.Println("keys:", sm.Keys())     // [b c d a]
	fmt.Println("values:", sm.Values()) // [10 20 20 30]

	sm.Put("c", 5) // value 改变 → 自动重新排序

	fmt.Println("after update keys:", sm.Keys()) // [c b d a]
}
