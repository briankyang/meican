package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ReadFromJsonFile(path string, config interface{}) {
	arr, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalln(err)
	}

	ReadFromJsonBytes(&arr, config)
}

//func init() {
//	log.Println("init config")
//}

func ReadFromJsonBytes(b *[]byte,config interface{}) {
	err := json.Unmarshal(*b, config)

	if err != nil {
		log.Fatal(err)
	}
}