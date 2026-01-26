package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {

	jsonStr := struct {
		Name       string    `json:"name,omitempty"`
		Age        int       `json:"age"`
		StartTieme time.Time `json:"start_time"`
	}{

		Age:        10,
		StartTieme: time.Now(),
	}
	jsonp, _ := json.Marshal(jsonStr)
	fmt.Println(string(jsonp))
}
