package main

import (
	"encoding/base64"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	base64Data := "jDO1d9RAMP0RcYDdJP/9jOT/Lv/+do0ACR12AAAAADD/CLndGBAEGJi+ByARKMgBMAM45b0HQAFIDlIZCJW+BxoJYm90MTIyNjQ1OMhZSANQAZABAVIbCLm9BxACGglib3QxMjI1NTM4uB1IA1ABkAEBUhwI5b0HEAMaCWJvdDEyMjU5NzjsgAJIA1ABkAEBUhkImL4HEAQaCWJvdDEyMjY0ODi4DVABkAEBUhkIo74HEAUaCWJvdDEyMjY1OTj"

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("%v", data)
}
