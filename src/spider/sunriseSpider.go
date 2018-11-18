package spider

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

func getHttpReq() {

}


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
	//req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	//	//req.Header.Set("Content-Type", "application/json")
	//	//client := &http.Client{}
	//	//
	//	//response, _ := client.Do(req)
	//	//defer response.Body.Close()
	//	//
	//    //body, _ := ioutil.ReadAll(response.Body)
	//    //fmt.Println(string(body))
	//	//jsonData, _ := simplejson.NewJson(body)
	jsonData, _ := GetJsonData(url, "POST", map[string]string{"Content-Type": "application/json"}, data)

	token, _ := 	jsonData.Get("data").Get("token").String()
  	return 	token
}

func MakeTokenHeader(r http.Request) http.Request{
	//def make_token_headers(self):
	//headers = self.headers.copy()
	//token = self.get_token()
	//headers["Authorization"] = "Bearer {}".format(token)
	//return headers
	token := GetToken()
	r.Header.Set("Authorization", "Bearer " + token)
	return r
}

func GetPcategorys(categoryChan chan string) {
	//def get_pcategorys(self):
	//url = "http://srmemberapp.srgow.com/goods/pcategorys/"
	//datas = self.get_data(url)
	//for data in datas:
	//yield data.get('id'), data.get('name')
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
}

func MakeCategoryPage(categoryChan chan string, urlChan chan string) {
	for {
		categoryId := <- categoryChan
		fmt.Println("is me make")
		for i := 1; i < 100; i ++ {
			url := fmt.Sprintf("http://srmemberapp.srgow.com/goods/search/%d?a=a&key=&category=%s", i, categoryId)
			fmt.Println(url)
			urlChan <- url
		}
	}

}

func GetOnePageGoods(urlChan chan string) {
	for {
		fmt.Println("还剩", len(urlChan), "个url")
		token := GetToken()
		headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
		datas, _ := GetJsonData(<- urlChan, "GET", headers, "")
		//abbid, _ := datas.Get("data").String()
		fmt.Println(datas.Get("data").GetIndex(0))
	}
}
