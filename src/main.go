package main

import (
	_ "net/http/pprof"
	"spider"
	"fmt"
	"time"
	"sync"
	"model"
	"domer"
	)

func getAll(){
	fmt.Println(time.Now())
	wg := new(sync.WaitGroup)
	categoryChan := make(chan string, 99)
	urlChan := make(chan [2]string, 100000)
	goodChan := make(chan model.Good, 100000)
	token := spider.GetToken()
	fmt.Println(token)
	wg.Add(1)
	go spider.GetPcategorys(categoryChan, wg)
	wg.Add(1)
	go spider.MakeCategoryPage(categoryChan, urlChan, wg)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go spider.GetOnePageGoods(urlChan, goodChan, token, wg)}
		wg.Add(1)
	go domer.DomToXlsx(goodChan, "good.xlsx", wg)
	wg.Wait()
	fmt.Println("exit")
	fmt.Println(time.Now())
}

func getGoodsInfo(abiids []string) (infos []model.GoodPriceInfo) {
	token := spider.GetToken()
	for _, abiid := range abiids{
		infos = append(infos, spider.GetAGood(abiid, token))
	}
	fmt.Println(infos)
	return
}

func main() {
	abiids := domer.ReadXlsx("/home/lzj/1.xlsx")
	getGoodsInfo(abiids)
}