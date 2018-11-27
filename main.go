package main

import (
	"flag"
	"log"
	"meican/config"
	"meican/service"
	"meican/util"
)

func main() {

	var users []service.User

	userConfigPath := flag.String("config", config.UserConfigPath, "user configuration path, ex: /Users/root/user.json")

	util.ReadFromJsonFile(*userConfigPath, &users)

	for _, u := range users {
		go u.Order()
	}

	for _ = range users {
		log.Println("complete with message: ", <-service.Complete)
	}
}
