package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"meican/model"
	"meican/util"
	"strings"
	"time"
)

var Complete chan interface{} = make(chan interface{})

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Exclusive []string `json:"exclusive"`
	Favorite []string `json:"favorite"`

	client util.MeiCanClient
}

func (p *User) Order() {

	if err := p.client.Login(p.Username, p.Password); err != nil {

		Complete <- errors.New(fmt.Sprintf("%s\t%s\t%s", err, p.Username, p.Password))
	}

	tabs, err := p.client.ListTab()

	if err != nil {
		Complete <-  err
		return
	}

	var tab *model.Tab
	if len(tabs) > 0 {
		tab = &tabs[0]
	}

	if tab == nil {
		Complete <-  errors.New("No tab available")
		return
	}

	restaurants, err := p.client.ListRestaurant(tab)

	if err != nil {
		Complete <-  err
		return
	}

	if restaurants == nil {
		Complete <-  errors.New("No Restaurant")
		return
	}

	//log.Println(len(restaurants), restaurants)

	var dishes []model.Dish
	for _, v := range restaurants  {
		d, _ := p.client.ListDishes(tab, &v)
		dishes = append(dishes, d...)
	}

	//log.Println(len(dishes), dishes)

	if len(dishes) <= 0 {
		Complete <-  errors.New("No dish offered")
		return
	}

	var targetDishes []model.Dish
	for _, v := range dishes {
		tag := true
		for _, exclusive := range p.Exclusive  {
			if strings.Contains(v.Name, exclusive) {
				tag = false
			}
		}
		if tag && v.Price > 0 {
			targetDishes = append(targetDishes, v)
		}
	}

	if len(targetDishes) <= 0 {
		Complete <-  errors.New("No target dish found, every dish is being excluded")
		return
	}

	var targetDish *model.Dish
	for _, v := range targetDishes{
		if targetDish != nil {
			break;
		}
		for _, favorite := range p.Favorite {
			if strings.Contains(v.Name, favorite) {
				log.Println("Find one favorite", v)
				targetDish = &v
				break
			}
		}
	}

	//log.Println(len(targetDishes), targetDishes)

	if targetDish == nil {
		rand.Seed(time.Now().Unix())
		idx := rand.Intn(len(targetDishes))
		targetDish = &targetDishes[idx]
	}

	log.Println("will order: ", targetDish.Name)

	res, _ := p.client.Order(tab, targetDish)

	Complete <- res
}