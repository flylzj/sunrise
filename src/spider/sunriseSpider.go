package spider

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	_ "github.com/mattn/go-sqlite3"
)


type Good struct {
	Abiid        int
	Mainname     string
	Subtitle     string
	Brandname    string
	Categoryname string
	Price        int
	Stock        string
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
	jsonData, _ := GetJsonData(url, "POST", map[string]string{"Content-Type": "application/json"}, data)
	token, _ := 	jsonData.Get("data").Get("token").String()
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

func GetOnePageGoods(urlChan chan [2]string, goodChan chan Good, token string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var urlArr [2]string
		select {
		    case urlArr = <- urlChan:
		    	//
		    case <- time.After(time.Second):
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
			brandname, _ := good.Get("brandname").String()
			categoryname, _ := good.Get("categoryname").String()
			price, _ := good.Get("price").Int()
			stock, _ := good.Get("stock").String()
			Gooda := Good{abiid, mainname, subtitle, brandname, categoryname, price, stock}
			goodChan <- Gooda
			//fmt.Println(Gooda)
		}
	}
}

func DomToFile(goodChan chan Good, filename string, wg *sync.WaitGroup) {
	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0755)
	for{
		//select {
		//case goodArr := <-:
		//
		//}
		var good Good
		sign := false
		select {
		    case good = <-goodChan:
		    	//
		    case <- time.After(time.Second * 5):
		    	fmt.Println("no goods")
		    	sign = true
		}
		if sign {
			break
		}
		//fmt.Println("还剩", len(goodChan), "个物品")
		f.Seek(0, 2)
		f.WriteString(strconv.Itoa(good.Abiid) + "," + good.Mainname + "," + strconv.Itoa(good.Price)+ "," + good.Stock + "\n")
	}
	defer wg.Done()
	defer f.Close()
    defer close(goodChan)
}

func DomToDB(goodChan chan Good, wg *sync.WaitGroup) {
	sqlStr := fmt.Sprintf("INSERT OR REPLACE INTO good(abiid, mainname, price, stock) values")
	for {
		fmt.Println("还剩", len(goodChan), "个物品")
		var good Good
		sign := false
		select {
		case good = <-goodChan:
			//
		case <-time.After(time.Second * 5):  //因为网络延迟需要等待一定的时间才能确定是没有商品了
			fmt.Println("no goods")
			sign = true
		}
		if sign {
			break
		}

		values := fmt.Sprintf("(%d, '%s', %d, '%s'),", good.Abiid, strings.Replace(good.Mainname, "'", "\"", -1), good.Price, good.Stock)
		sqlStr += values
	}
	sqlStr = strings.Trim(sqlStr, ",") // 删掉sql末尾的逗号
	db, _ := sql.Open("sqlite3", "good.db")
	res, err := db.Exec(sqlStr)
	if err != nil{
		fmt.Println(err)
	}
	res.RowsAffected()
	defer db.Close()
	defer wg.Done()
}

func InitDB() *sql.DB{
	db, err := sql.Open("sqlite3", "good.db")
	if err != nil {
		fmt.Println("err", err)
	}
	db.Exec("DROP TABLE IF EXISTS good")
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS good (abiid INT PRIMARY KEY NOT NULL, mainname VARCHAR(64), price INT NOT NULL , stock VARCHAR(12) NOT NULL)")
	if err != nil {
		fmt.Println("err", err)
	}
	statement.Exec()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("init db success")
	return db
}

