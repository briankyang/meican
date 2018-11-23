package model

type Tab struct {
	Title string
	TargetTime interface{}
	Status TabStatus
	Uid string
	Address []Address
}

const (
	Available  = TabStatus("AVAILABLE")
	Closed  = TabStatus("CLOSED")
	NotYet  = TabStatus("NOT_YET")
	Unknown = TabStatus("UNKNOWN")
	Ordered = TabStatus("ORDERED")
)

var statusMap = map[string]TabStatus {
	string(Available): Available,
	string(Closed): Closed,
	string(NotYet): NotYet,
	string(Unknown): Unknown,
	string(Ordered): Ordered,
}

type TabStatus string

func (ts TabStatus) Status() string  {
	return string(ts)
}

func ToTabStatus(status string) TabStatus {
	return statusMap[status]
}

func NewTab(p interface{}) *Tab {

	t, _ := p.(map[string]interface{})

	l, _ := t["calendarItemList"].([]interface{})

	if len(l) <= 0 {
		return nil
	}

	overView := l[0].(map[string]interface{})
	uInfo := overView["userTab"].(map[string]interface{})
	uAddress := uInfo["corp"].(map[string]interface{})["addressList"].([]interface{})

	var addrs []Address

	for _, addr := range uAddress {
		addrs = append(addrs, NewAddress(&addr))
	}

	tab := &Tab{
		Title: overView["title"].(string),
		Status: ToTabStatus(overView["status"].(string)),
		TargetTime: overView["targetTime"],
		Uid: uInfo["uniqueId"].(string),
		Address: addrs,
	}

	return tab
}