package main

import (
	"log"
	"meican/config"
	"meican/service"
	"meican/util"
)


func main()  {

	var users []service.User

	util.ReadFromJsonFile(config.UserConfigPath , &users)

	for _, u := range users{
		log.Println("1231231")
		go u.Order()
	}

	for _ = range users {
		log.Println("complete with message: ", <- service.Complete)
	}
}