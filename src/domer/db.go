package domer

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"model"
	"strings"
	"sync"
	"time"
)

func DomToDB(goodChan chan model.Good, wg *sync.WaitGroup) {
	sqlStr := fmt.Sprintf("INSERT OR REPLACE INTO good(abiid, mainname, price, stock) values")
	for {
		fmt.Println("还剩", len(goodChan), "个物品")
		var good model.Good
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

func CreateTmpStockTable(){
	db, err := sql.Open("sqlite3", "good.db")
	if err != nil {
		fmt.Println("err", err)
	}
	db.Exec("DROP TABLE IF EXISTS good")
	db.Exec("CREATE TABLE tmp_stock(abiid INT PRIMARY KEY NOT NULL, stock_num INT NOT NULL )")
}

func DomTmpStock(info model.GoodPriceInfo){
	db, err := sql.Open("sqlite3", "good.db")
	defer db.Close()
	if err != nil {
		fmt.Println("err", err)
	}
	stock_num := SearchTmpStock(info)
	if stock_num == -1 {
		_, err = db.Exec(fmt.Sprintf("INSERT INTO tmp_stock(abiid, stock_num) VALUES(%s, %d)", info.Abiid, info.Num))
		if err != nil {
			fmt.Println(err)
			fmt.Println("插入", info.Abiid, "失败")
			return
		}
		fmt.Println("插入", info.Abiid,"成功")
	}
	return
}

func SearchTmpStock(info model.GoodPriceInfo) int{
	db, err := sql.Open("sqlite3", "good.db")
	if err != nil {
		fmt.Println("err", err)
	}
	rows, err := db.Query(fmt.Sprintf("SELECT abiid, stock_num FROM tmp_stock WHERE abiid=%s", info.Abiid))
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	var abiid int
	var stock_num int
	for rows.Next(){
		err = rows.Scan(&abiid, &stock_num)
	}
	if abiid == 0{
		return -1
	}
	return stock_num
}