package util

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	errors2 "meican/errors"
	"meican/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	format        = "2006-01-02"
	format_normal = "2006-01-02 15:04:05"
)

type MeiCanConfig struct {
	BaseUrl string `json:"base_url"`

	LoginPath      string `json:"login_path"`
	CalendarPath   string `json:"calendar_path"`
	RestaurantPath string `json:"restaurant_path"`
	DishPath       string `json:"dish_path"`
	OrderPath      string `json:"order_path"`
}

type MeiCanClient struct {
	config *MeiCanConfig

	httpClient *http.Client
}

func (p *MeiCanClient) assure() {
	if p.config == nil {
		p.config = &defaultConfig
	}

	if p.httpClient == nil {
		p.httpClient = http.DefaultClient
		p.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return errors2.StopRedirectError
		}
	}

	if p.httpClient.Jar == nil {

		jar, err := cookiejar.New(nil)

		if err != nil {
			log.Fatalln("cannot instance a cookieJar")
		}

		p.httpClient.Jar = jar
	}
}

func (p *MeiCanClient) post(url string, contentType string, content string) (string, error) {
	p.assure()

	content = strings.Replace(content, "+", "%20", -1)

	//log.Println(content)

	return p.request("post", url, contentType, content)
}

func (p *MeiCanClient) get(url string) (string, error) {
	p.assure()

	return p.request(http.MethodGet, url, "", "")
}

func (p *MeiCanClient) request(method string, endpoint string, contentType string, content string) (string, error) {

	var contentReader io.Reader
	if len(content) == 0 {
		contentReader = nil
	} else {
		contentReader = strings.NewReader(content)
	}

	req, _ := http.NewRequest(method, endpoint, contentReader)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36")

	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := p.httpClient.Do(req)

	if err != nil && !strings.HasSuffix(err.Error(), errors2.StopRedirectError.Error()) {
		return "", err
	}

	contentBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "PLAY_FLASH" {
			hint, _ := url.ParseQuery(cookie.Value)

			errorHint := hint.Get("error")

			if len(errorHint) != 0 {
				return "", errors.New(errorHint)
			}
		}
	}

	return string(contentBytes), nil
}

func (p *MeiCanClient) Login(username string, password string) error {

	p.assure()

	endpoint := p.config.BaseUrl + p.config.LoginPath

	content := url.Values{"username": {username}, "password": {password}, "loginType": {"username"}, "remember": {"true"}}

	//log.Println("start login", endpoint, content.Encode())

	_, e := p.post(endpoint, "application/x-www-form-urlencoded", content.Encode())

	return e
}

func (p *MeiCanClient) ListTab() ([]model.Tab, error) {
	p.assure()

	now := time.Now()

	content := url.Values{"beginDate": {now.Format(format)}, "endDate": {now.Add(time.Hour * 7 * 24).Format(format)}, "withOrderDetail": {"false"}}

	endpoint := p.config.BaseUrl + p.config.CalendarPath + "?" + content.Encode()

	r, e := p.get(endpoint)

	if e != nil {
		return nil, e
	}

	var cc interface{}

	json.Unmarshal([]byte(r), &cc)

	m, _ := cc.(map[string]interface{})

	dateList, _ := m["dateList"].([]interface{})
	var tabs []model.Tab

	for _, v := range dateList {
		t := model.NewTab(v)

		if t != nil && t.Status == model.Available {
			tabs = append(tabs, *t)
		}
	}

	return tabs, nil
}

func (p *MeiCanClient) ListRestaurant(t *model.Tab) ([]model.Restaurant, error) {
	p.assure()

	content := url.Values{"tabUniqueId": {t.Uid}, "targetTime": {string(time.Unix(int64(t.TargetTime.(float64)/1000), 0).Format(format_normal))}}

	endpoint := p.config.BaseUrl + p.config.RestaurantPath + "?" + content.Encode()

	//log.Println(endpoint)

	r, e := p.get(endpoint)

	if e != nil {
		return nil, e
	}

	var cc interface{}

	json.Unmarshal([]byte(r), &cc)

	m, _ := cc.(map[string]interface{})

	restaurantList, _ := m["restaurantList"].([]interface{})

	var restaurants []model.Restaurant

	for _, v := range restaurantList {
		r := model.NewRestaurant(&v)
		restaurants = append(restaurants, *r)
	}

	return restaurants, nil
}

func (p *MeiCanClient) ListDishes(t *model.Tab, r *model.Restaurant) ([]model.Dish, error) {

	endpoint := p.config.BaseUrl + p.config.DishPath

	content := url.Values{"tabUniqueId": {t.Uid}, "targetTime": {string(time.Unix(int64(t.TargetTime.(float64)/1000), 0).Format(format_normal))}, "restaurantUniqueId": {r.UniqueId}}

	res, e := p.post(endpoint, "application/x-www-form-urlencoded", content.Encode())

	if e != nil {
		return nil, e
	}

	var cc interface{}

	json.Unmarshal([]byte(res), &cc)

	m, _ := cc.(map[string]interface{})

	dishList, _ := m["dishList"].([]interface{})

	var dishes []model.Dish

	for _, v := range dishList {
		d := model.NewDish(&v)
		dishes = append(dishes, *d)
	}
	return dishes, nil
}

func (p *MeiCanClient) Order(t *model.Tab, d *model.Dish) (string, error) {

	order := "[{\"dishId\":" + strconv.Itoa(d.Id) + ",\"count\":1" + "}]"
	remark := "[{\"dishId\":\"" + strconv.Itoa(d.Id) + "\",\"remark\":\"\"}]"

	content := url.Values{"tabUniqueId": {t.Uid},
		"order":               {order},
		"targetTime":          {string(time.Unix(int64(t.TargetTime.(float64)/1000), 0).Format(format_normal))},
		"userAddressUniqueId": {t.Address[0].Uuid},
		"corpAddressUniqueId": {t.Address[0].Uuid},
		"corpAddressRemark":   {""},
		"remarks":             {remark},
	}

	return p.post(p.config.BaseUrl+p.config.OrderPath, "application/x-www-form-urlencoded", content.Encode())
}

var defaultConfig MeiCanConfig

func init() {
	defaultConfig = MeiCanConfig{
		BaseUrl: "https://meican.com",

		LoginPath:      "/account/directlogin",
		CalendarPath:   "/preorder/api/v2.1/calendarItems/list",
		RestaurantPath: "/preorder/api/v2.1/restaurants/list",
		DishPath:       "/preorder/api/v2.1/restaurants/show",
		OrderPath:      "/preorder/api/v2.1/orders/add",
	}
	//log.Println("init meican")
	//ReadFromJsonFile(config.UrlConfigPath, &defaultConfig)
}
