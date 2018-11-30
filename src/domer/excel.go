package domer

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"model"
	"strconv"
	"sync"
	"time"
)

func DomToXlsx(goodChan chan model.Good, filename string, wg *sync.WaitGroup) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//Abiid        int
	//Mainname     string
	//Subtitle     string
	//Brandid      int
	//Brandname    string
	//CategoryId   int
	//Categoryname string
	//Price        int
	//Stock        string
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "abiid"
	cell = row.AddCell()
	cell.Value = "商品名"
	cell = row.AddCell()
	cell.Value = "subtitle"
	cell = row.AddCell()
	cell.Value = "brandid"
	cell = row.AddCell()
	cell.Value = "brandname"
	cell = row.AddCell()
	cell.Value = "categoryid"
	cell = row.AddCell()
	cell.Value = "categoryname"
	cell = row.AddCell()
	cell.Value = "price"
	cell = row.AddCell()
	cell.Value = "stock"
	cell = row.AddCell()
	cell.Value = "intstock"
	for{
		var good model.Good
		sign := false
		select {
		case good = <-goodChan:
			//
		case <- time.After(time.Second * 10):
			fmt.Println("no goods")
			sign = true
		}
		if sign {
			break
		}
		fmt.Println("还剩", len(goodChan), "个物品")
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = strconv.Itoa(good.Abiid)
		cell = row.AddCell()
		cell.Value = good.Mainname
		cell = row.AddCell()
		cell.Value = good.Subtitle
		cell = row.AddCell()
		cell.Value = good.Brandid
		cell = row.AddCell()
		cell.Value = good.Brandname
		cell = row.AddCell()
		cell.Value = good.CategoryId
		cell = row.AddCell()
		cell.Value = good.Categoryname
		cell = row.AddCell()
		cell.Value = strconv.Itoa(good.Price)
		cell = row.AddCell()
		cell.Value = good.Stock
		cell = row.AddCell()
		cell.Value = strconv.Itoa(good.IntStock)
	}
	err = file.Save(filename)
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer wg.Done()
	//defer close(goodChan)
}

func ReadXlsx(filename string){

}
