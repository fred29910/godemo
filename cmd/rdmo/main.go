package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redisaabbcc", // no password set
		DB:       0,             // use default DB
	})

	// 添加一些成员到集合
	err := rdb.SAdd(ctx, "myset", "member1", "member2", "member3").Err()
	if err != nil {
		fmt.Println("Error adding members to set:", err)
		return
	}

	// 检查成员是否存在于集合中
	isMember, err := rdb.SIsMember(ctx, "myset", "member2").Result()
	if err != nil {
		fmt.Println("Error checking if member exists:", err)
		return
	}

	if isMember {
		fmt.Println("member2 exists in myset")
	} else {
		fmt.Println("member2 does not exist in myset")
	}

	// 也可以检查一个不存在的成员
	isMember, err = rdb.SIsMember(ctx, "myset", "member4").Result()
	if err != nil {
		fmt.Println("Error checking if member exists:", err)
		return
	}

	if isMember {
		fmt.Println("member4 exists in myset")
	} else {
		fmt.Println("member4 does not exist in myset")
	}

	// 也可以检查一个不存在key
	isMember, err = rdb.SIsMember(ctx, "myset1xiaimsdiaimsma", "member4").Result()
	if err != nil {
		fmt.Println("Error checking if member exists:", err)
		return
	}
}

func main() {
	ExampleClient()
}
