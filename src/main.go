package main

import (
	"domer"
	"fmt"
	"model"
	"noticer"
	"os"
	"spider"
	"sync"
	"time"
	"tk"
)

func getAll(root *tk.Root){
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
		go spider.GetOnePageGoods(urlChan, goodChan, token, wg)
	}
	wg.Add(1)
	go domer.DomToXlsx(goodChan, "good.xlsx", wg, root)
	wg.Wait()
	fmt.Println("exit")
	fmt.Println(time.Now())
}

func getGoodsInfo(filename string, root *tk.Root){
	domer.CreateTmpStockTable()
	abiids := domer.ReadXlsx(filename)
	infoChan := make(chan model.GoodPriceInfo)
	token := spider.GetToken()
	for _, abiid := range abiids{
		go spider.GetAGood(abiid, token, infoChan)
	}
	for {
		timeout := false
		select {
		case info := <- infoChan:
			domer.DomTmpStock(info, root)
		case <- time.After(time.Second * 2):
		    timeout = true
		}
		if timeout{
			break
		}
	}
	infos := domer.SearchNeedNotice()
	if len(infos) != 0 {
		domer.DomStockToExcel("output.xlsx", infos)
		noticer.SendEmail()
		root.MessageChan <- [2]string{"main", "邮件发送成功"}
	}else {
		root.MessageChan <- [2]string{"main", "商品库存没有变化"}
	}
}

func startMonitor(root *tk.Root){
	for{
		hour := root.Spider.HourSpinButton.GetValueAsInt()
		minute := root.Spider.MinuteSpinButton.GetValueAsInt()
		filename := root.Exceler.ExcelEntry.GetText()
		_, err := os.Open(filename)
		println(err)
		if os.IsNotExist(err){
			root.MessageChan <- [2]string{"main", "请选择文件"}
			continue
		}
		root.MessageChan <- [2]string{"main", "开始爬取"}
		getGoodsInfo(filename, root)
		fmt.Println()
		root.MessageChan <- [2]string{"main", "爬取完成, 下次爬取将在："+time.Now().Add(time.Hour * time.Duration(hour) + time.Minute * time.Duration(minute)).String()}
		time.Sleep(time.Hour * time.Duration(hour) + time.Minute * time.Duration(minute))
	}
}

func bindStartButton(root *tk.Root){
	root.Spider.StartButton.Clicked(func() {
		go startMonitor(root)
		root.Spider.StartButton.SetSensitive(false)
	})
}

func bindOutputButton(root *tk.Root){
	root.OutputSpider.OutputButton.Clicked(func() {
		go getAll(root)
	})
}

func main() {
	token := spider.GetToken()
	fmt.Println(token)
	spider.GetAGood("236081", token, nil)
	//r := tk.MakeRootWindow()
	//r.Window.ShowAll()
	//bindStartButton(r)
	//bindOutputButton(r)
	//gtk.Main()
}