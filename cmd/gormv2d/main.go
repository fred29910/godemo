package main

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"time"
)

func main() {

	buffer := make([]byte, 1024)

	leftByte := 0
	data := []byte(randomString(129))
	leftByte += len(data)
	copy(buffer, data)

	processedLength := leftByte

	copy(buffer, buffer[processedLength:leftByte])

	leftByte -= processedLength

	data2 := []byte(randomString(100))

	bufReader := bytes.NewBuffer(data2)

	l, err := bufReader.Read(buffer[leftByte:])
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("read %d bytes", l)

	fmt.Println(string(buffer))

}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // 设定随机种子
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
