package domer

import (
	"database/sql"
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
