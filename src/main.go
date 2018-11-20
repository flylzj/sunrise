package main

import (
	"net/http"
	"spider"
	"sync"
	_ "net/http/pprof"
)

func main() {
	var wg sync.WaitGroup
	categoriyChan := make(chan string, 99)
	urlChan := make(chan [2]string, 10000)
	goodChan := make(chan spider.Good, 10000)
	wg.Add(12000)
	token := spider.GetToken()
	go spider.GetPcategorys(categoriyChan)
    go spider.MakeCategoryPage(categoriyChan, urlChan)
	for i := 0; i < 1000; i++ {
		go spider.GetOnePageGoods(urlChan, goodChan, token)
	}
	go spider.DomToFile(goodChan, "good.txt")
	go func(){
		http.ListenAndServe(":6060",nil)
	}()
	wg.Wait()
}