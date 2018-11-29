package main

import (
	"flag"
	"log"
	"meican/config"
	"meican/service"
	"meican/util"
)

var userConfigPath string

func init() {
	flag.StringVar(&userConfigPath, "config", config.UserConfigPath, "user configuration path, ex: /Users/root/user.json")
}

func main() {

	flag.Parse()

	var users []service.User

	util.ReadFromJsonFile(userConfigPath, &users)

	for idx, _ := range users {
		go users[idx].Order()
	}

	for _ = range users {
		log.Println("complete with message: ", <-service.Complete)
	}
}
