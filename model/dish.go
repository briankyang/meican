package model

type Dish struct {
	Id int
	Name string
	Price float64
}

func NewDish(p *interface{}) *Dish {
	d, _ := (*p).(map[string]interface{})

	r := &Dish{
		Id: int(d["id"].(float64)),
		Name: d["name"].(string),
		Price: d["originalPriceInCent"].(float64),
	}

	return r
}