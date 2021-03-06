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
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Exclusive []string `json:"exclusive"`
	Favorite  []string `json:"favorite"`

	client util.MeiCanClient
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (p *User) Order() (string, error) {

	if err := p.client.Login(p.Username, p.Password); err != nil {

		return "", errors.New(fmt.Sprintf("%s\t%s\t%s", err, p.Username, p.Password))
	}

	log.Println(fmt.Sprintf("[%s]登陆成功, 开始订餐～", p.Username))

	tabs, err := p.client.ListTab()

	if err != nil {
		return "", err
	}

	var tab *model.Tab
	if len(tabs) > 0 {
		tab = &tabs[0]
	}

	if tab == nil {
		return "", errors.New(fmt.Sprintf("[%s]No tab available", p.Username))
	}

	restaurants, err := p.client.ListRestaurant(tab)

	if err != nil {
		return "", err
	}

	if restaurants == nil {
		return "", errors.New(fmt.Sprintf("[%s]No Restaurant", p.Username))
	}

	//log.Println(len(restaurants), restaurants)

	var dishes []model.Dish
	for _, v := range restaurants {
		d, _ := p.client.ListDishes(tab, &v)
		dishes = append(dishes, d...)
	}

	//log.Println(len(dishes), dishes)

	if len(dishes) <= 0 {
		return "", errors.New(fmt.Sprintf("[%s]No dish offered", p.Username))
	}

	var targetDishes []model.Dish
	for _, v := range dishes {
		tag := true
		for _, exclusive := range p.Exclusive {
			if strings.Contains(v.Name, exclusive) {
				tag = false
			}
		}
		if tag && v.Price > 0 {
			targetDishes = append(targetDishes, v)
		}
	}

	if len(targetDishes) <= 0 {
		return "", errors.New(fmt.Sprintf("[%s]No target dish found, every dish is being excluded", p.Username))
	}

	var favoriteDishes []model.Dish
	for idx, v := range targetDishes {
		for _, favorite := range p.Favorite {
			if strings.Contains(targetDishes[idx].Name, favorite) {
				favoriteDishes = append(favoriteDishes, v)
			}
		}
	}

	var targetDish *model.Dish
	if len(favoriteDishes) == 0 {
		idx := rand.Intn(len(targetDishes))
		targetDish = &targetDishes[idx]
		log.Println("choose dish by random", targetDish)
	} else {
		idx := rand.Intn(len(favoriteDishes))
		targetDish = &favoriteDishes[idx]
		log.Println("choose favorite dish by random", targetDish)
	}

	log.Println(fmt.Sprintf("[%s]will order: [%s]", p.Username, targetDish.Name))

	res, _ := p.client.Order(tab, targetDish)

	return res, nil
}
