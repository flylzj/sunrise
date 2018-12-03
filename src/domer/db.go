package domer

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"model"
	"strconv"
	"strings"
	"sync"
	"time"
	"tk"
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
	//db.Exec("DROP TABLE IF EXISTS tmp_stock")
	_, err = db.Exec("CREATE TABLE tmp_stock(abiid INT PRIMARY KEY NOT NULL, mainname VARCHAR(255),price INT NOT NULL,realprice INT, stock VARCHAR(255) NOT NULL, stock_num INT NOT NULL, stock_num_2 INT NOT NULL)")
}

func DomTmpStock(info model.GoodPriceInfo, root *tk.Root) int{
	db, err := sql.Open("sqlite3", "good.db")
	defer db.Close()
	if err != nil {
		fmt.Println("err", err)
	}
	stockNum := SearchTmpStock(info)
	if stockNum == -1 {
		_, err = db.Exec(fmt.Sprintf("INSERT INTO tmp_stock(abiid, mainname, price, realprice, stock, stock_num, stock_num_2) VALUES(%s, '%s', %d, %d, '%s', %d, %d)", info.Abiid, info.Mainname, info.Price, info.RealPrice, info.Stock, info.Num, info.Num))
		if err != nil {
			fmt.Println(err)
			fmt.Println("插入", info.Abiid, "失败")
		}else {
			fmt.Println("插入", info.Abiid,"成功")
		}
		root.MessageChan <- [2]string{"main", "插入"+info.Abiid+"失败"}
	}else {
		_, err = db.Exec(fmt.Sprintf("UPDATE tmp_stock SET stock='%s',price=%d, price=%d, stock_num=%d, stock_num_2=%d WHERE abiid=%s", info.Stock, info.Price, info.RealPrice, stockNum, info.Num, info.Abiid))
		if err != nil {
			fmt.Println(err)
			fmt.Println("更新", info.Abiid, "失败")
		}else {
			fmt.Println("更新", info.Abiid,"成功")
		}
		root.MessageChan <- [2]string{"main", "更新"+info.Abiid+"成功"}
	}
	return stockNum
}

func SearchTmpStock(info model.GoodPriceInfo) int{
	db, err := sql.Open("sqlite3", "good.db")
	if err != nil {
		fmt.Println("err", err)
	}
	rows, err := db.Query(fmt.Sprintf("SELECT abiid, stock_num_2 FROM tmp_stock WHERE abiid=%s", info.Abiid))
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	var abiid int
	var stockNum int
	for rows.Next(){
		err = rows.Scan(&abiid, &stockNum)
	}
	if abiid == 0{ // 数据库里没有该商品
		return -1
	}
	return stockNum
}

func SearchNeedNotice()(infos []model.GoodPriceInfo){
	db, err := sql.Open("sqlite3", "good.db")
	if err != nil {
		fmt.Println("err", err)
	}
	rows, err := db.Query("SELECT abiid, mainname, price, realprice, stock, stock_num, stock_num_2 FROM tmp_stock WHERE stock_num != stock_num_2")
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next(){
		var (
			abiid int
			mainname string
			price int
			realprice int
			stock string
			stock_num int
			stock_num_2 int
		)
		err = rows.Scan(&abiid, &mainname, &price, &realprice, &stock, &stock_num, &stock_num_2)
		infos = append(infos, model.GoodPriceInfo{
			Abiid:strconv.Itoa(abiid),
			Mainname:mainname,
			Price:price,
			RealPrice:realprice,
			Stock:stock,
			Num:stock_num,
			Num2:stock_num_2})
	}
	return infos
}