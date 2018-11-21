package main

import (
	_ "net/http/pprof"
	"spider"
)

func main() {
	//wg := new(sync.WaitGroup)
	//categoryChan := make(chan string, 99)
	//urlChan := make(chan [2]string, 100000)
	//goodChan := make(chan spider.Good, 100000)
	//token := spider.GetToken()
	//wg.Add(1)
	//go spider.GetPcategorys(categoryChan, wg)
	//wg.Add(1)
    //go spider.MakeCategoryPage(categoryChan, urlChan,wg)
	//for i := 0; i < 100; i++ {
	//	wg.Add(1)
	//	go spider.GetOnePageGoods(urlChan, goodChan, token, wg)
	//}
	//wg.Add(1)
	//go spider.DomToFile(goodChan, "good.csv", wg)
	////wg.Add(100)
	////go func(){
	////	http.ListenAndServe(":6060",nil)
	////}()
	//wg.Wait()
	//fmt.Println("exit")
	spider.InitDB()
}