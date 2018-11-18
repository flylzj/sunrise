package main

import (
	"spider"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	categoriyChan := make(chan string, 99)
	urlChan := make(chan string, 10000)
	wg.Add(2000)
	go spider.GetPcategorys(categoriyChan)
    go spider.MakeCategoryPage(categoriyChan, urlChan)
	for i := 0; i < 1000; i++ {
		go spider.GetOnePageGoods(urlChan)
	}
	wg.Wait()
	//for ca := range categoriyChan {
	//	fmt.Println(ca)
	//}
	//close(categoriyChan)
}