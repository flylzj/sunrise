package domer

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"model"
	"os"
	"strconv"
	"sync"
	"time"
	"tk"
)

func DomToXlsx(goodChan chan model.Good, filename string, wg *sync.WaitGroup, root *tk.Root) {
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
			//root.MessageChan <- [2]string{"output", "保存"+strconv.Itoa(good.Abiid)+"成功"}
		case <- time.After(time.Second * 10):
			root.MessageChan <- [2]string{"output", "导出完成"}
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

func ReadXlsx(filename string) (abiids []string){
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	sheet := xlFile.Sheets[0]
	for _, row := range sheet.Rows{
		abiid := row.Cells[0]
		abiids = append(abiids, abiid.Value)
	}
	return

}

func DomStockToExcel(filename string, infos []model.GoodPriceInfo){
	os.Remove(filename)
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row := sheet.AddRow()
	row.AddCell().Value = "abiid"
	row.AddCell().Value = "mainname"
	row.AddCell().Value = "price"
	row.AddCell().Value = "stock"
	row.AddCell().Value = "num1"
	row.AddCell().Value = "num2"
	for _, info := range infos{
		row = sheet.AddRow()
		row.AddCell().Value = info.Abiid
		row.AddCell().Value = info.Mainname
		row.AddCell().Value = strconv.Itoa(info.Price)
		row.AddCell().Value = info.Stock
		row.AddCell().Value = strconv.Itoa(info.Num)
		row.AddCell().Value = strconv.Itoa(info.Num2)
	}
	err = file.Save(filename)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
