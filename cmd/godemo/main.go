package main

import (
	"fmt"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	viper.AddRemoteProvider("consul", "localhost:8500", "app/demo1")
	viper.SetConfigType("json") // Need to explicitly set this to json
	err := viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(viper.Get("port"))     // 8080
	fmt.Println(viper.Get("hostname")) // myhostname.com
}
