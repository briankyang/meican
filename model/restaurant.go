package model

type Restaurant struct {
	UniqueId string
	Name string
	Open bool
	Rating int
	Tel string
	Latitude float64
	Longitude float64
}

func NewRestaurant(p *interface{}) *Restaurant {
	m, _ := (*p).(map[string]interface{})

	r := &Restaurant{
		UniqueId: m["uniqueId"].(string),
		Name: m["name"].(string),
		Open: m["open"].(bool),
		Rating: int(m["rating"].(float64)),
		Tel: m["tel"].(string),
		Latitude:m["latitude"].(float64),
		Longitude: m["longitude"].(float64),
	}

	return r
}