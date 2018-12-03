package spider

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"model"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetToken() string{
	// python的函数，改成go
	//url = "http://srmemberapp.srgow.com/sys/token"
	//_nonce = str(randint(1001, 10000))
	//_timestamp = str(int(time.time()))
	//_array = [_nonce, self._appsecret, _timestamp]
	//_array.sort()
	//_tmp = ''.join(_array)
	//m = hashlib.md5(_tmp.encode())
	//_signature = m.hexdigest().upper()
	//data = {
	//	"appid": self._appid,
	//		"appsecret": self._appsecret,
	//		"timestamp": _timestamp,
	//		"signature": _signature,
	//		"nonce": _nonce
	//}
	//r = requests.post(url, headers=self.headers, data=data)
	//d = r.json()
	//token = d.get("data").get("token")
	//return token
	url := "http://srmemberapp.srgow.com/sys/token"
	appsecret := "e1d0b361201e4324b37c968fb71f0d3c"
	appid := "sunrise_member"
	nonce := fmt.Sprintf("%d", rand.Intn(9000) + 1001)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	_array := []string{appsecret, nonce, timestamp}
	sort.Strings(_array)
	_tmp := strings.Join(_array, "")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(_tmp))
	cipherStr := md5Ctx.Sum(nil)
	signature := hex.EncodeToString(cipherStr)
	data := fmt.Sprintf("{\"appid\": \"%s\", \"appsecret\": \"%s\", \"timestamp\": \"%s\", \"signature\": \"%s\", \"nonce\": \"%s\"}",
		appid, appsecret, timestamp, strings.ToUpper(signature), nonce)
	jsonData, _ := GetJsonData(url, "POST", map[string]string{"Content-Type": "application/json"}, data)
	token, _ := jsonData.Get("data").Get("token").String()
  	return 	token
}

func GetPcategorys(categoryChan chan string, wg *sync.WaitGroup) {
	//def get_pcategorys(self):
	//url = "http://srmemberapp.srgow.com/goods/pcategorys/"
	//datas = self.get_data(url)
	//for data in datas:
	//yield data.get('id'), data.get('name')
	defer wg.Done()
	url := "http://srmemberapp.srgow.com/goods/pcategorys/"
	token := GetToken()
	headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
	datas, _ := GetJsonData(url, "GET", headers, "")
	itemsCount, _ := datas.Get("data").Array()
	for i := 0; i < len(itemsCount); i ++ {
		item := datas.Get("data").GetIndex(i)
		categoryId, err := item.Get("id").String()
		if err != nil {
			continue
		}
		categoryChan <- categoryId
	}
	defer close(categoryChan)
}

func MakeCategoryPage(categoryChan chan string, urlChan chan [2]string, wg *sync.WaitGroup) {
	for{
		var categoryId string
		select {
		    case categoryId = <- categoryChan:
		    	//
		    case <- time.After(time.Second):
		    	categoryId = ""
		}
		if categoryId == ""{
			break
		}
		urlArr := [2]string{"1", categoryId}
		urlChan <- urlArr
	}
	defer wg.Done()
}

func makeUrl(urlArr [2]string) string{
	return fmt.Sprintf("http://srmemberapp.srgow.com/goods/search/%s?a=a&key=&category=%s",urlArr[0], urlArr[1])
}

func GetOnePageGoods(urlChan chan [2]string, goodChan chan model.Good, token string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var urlArr [2]string
		select {
		case urlArr = <- urlChan:
			//
		case <- time.After(time.Second * 5):
			fmt.Println("no url")
			urlArr[0] = ""
		}
		if urlArr[0] == "" {
			break
		}
		// urlArr := <- urlChan
		//fmt.Println("还剩", len(urlChan), "个url")
		headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
		datas, _ := GetJsonData(makeUrl(urlArr), "GET", headers, "")
		goods, _ := datas.Get("data").Array()
		goodsCount := len(goods)
		if goodsCount == 20{
			page, _ := strconv.Atoi(urlArr[0])
			if page % 10 == 1 {
				fmt.Println("page", page)
			    for i := 0; i < 10; i ++ {
			    	page = page + 1
					urlChan <- [2]string{strconv.Itoa(page + 1), urlArr[1]}
				}
			}
		}
		for i := 0; i < goodsCount; i ++ {
			good := datas.Get("data").GetIndex(i)
			mainname, _ := good.Get("mainname").String()
			abiid, _ := good.Get("abiid").Int()
			subtitle, _ := good.Get("subtitle").String()
			brandid, _ := good.Get("brandid").String()
			brandname, _ := good.Get("brandname").String()
			categoryid, _ := good.Get("categoryid").String()
			categoryname, _ := good.Get("categoryname").String()
			price, _ := good.Get("price").Int()
			stock, _ := good.Get("stock").String()
			intStock, _ := good.Get("intstock").Int()
			Gooda := model.Good{abiid, mainname, subtitle, brandid, brandname, categoryid, categoryname, price, stock, intStock}
			goodChan <- Gooda
		}
	}
}

func GetAGood(abiid string, token string, infoChan chan model.GoodPriceInfo){
	info := model.GoodPriceInfo{Abiid:abiid}
	GetGoodInfo(abiid, token, &info)
	GetGoodPrice(abiid, token, &info)
	fmt.Println(info)
	infoChan <- info
}

func GetGoodPrice(abiid string, token string, info *model.GoodPriceInfo){
	url := fmt.Sprintf("http://srmemberapp.srgow.com/goods/prices/%s", abiid)
	headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
	data, _ := GetJsonData(url, "GET", headers, "")
	data = data.Get("data")
	realprice, _ := data.Get("realprice").Int()
	price, _ := data.Get("price").Int()
	stock, _ := data.Get("stock").String()
	num, _ := data.Get("num").Int()
	info.RealPrice = realprice
	info.Price = price
	info.Stock = stock
	info.Num = num
}

func GetGoodInfo(abiid string, token string, info *model.GoodPriceInfo){
	url := fmt.Sprintf("http://b2carticleinfo.lib.cdn.srgow.com/api/v1/Article?languageid=1&abiid=%s", abiid)
	headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
	data, _ := GetJsonData(url, "GET", headers, "")
	mainname, _ := data.Get("mainname").String()
	info.Mainname = mainname
}

