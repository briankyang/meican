package model

type Address struct {
	Uuid string
	Address string
	PickUpLocation string
}

func NewAddress(p *interface{}) Address {
	addr := (*p).(map[string]interface{})

	return Address{
		Uuid: addr["uniqueId"].(string),
		Address: addr["address"].(string),
		PickUpLocation: addr["pickUpLocation"].(string),
	}
}