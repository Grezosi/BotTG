package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Good struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//создаем гет-параметр
		get := r.URL.Query().Get("id")
		// создаем подключение к редису
		conn, err := redis.Dial("tcp", "redis:6379")
		if err != nil {
			fmt.Println("Ошибка Редис", err)
		}
		intGet, _ := strconv.Atoi(get)
		cashData, _ := redis.String(conn.Do("GET", get))
		var db *sql.DB

		if cashData == "" {
			//создаем кеш
			if intGet >= 1 && intGet <= 10 {
				fmt.Println("Zashli v block s BD1")
				db, err = sql.Open("mysql", "root:test123@tcp(mysql1:3306)/mysql1")
				if err != nil {
					fmt.Println("Net podklucheniya k DB1", err)
				}
			} else if intGet >= 11 && intGet <= 20 {
				db, err = sql.Open("mysql", "root:test123@tcp(mysql2:3306)/mysql2")
				if err != nil {
					fmt.Println("Net podklucheniya k DB2", err)
				}

			} else if intGet >= 21 && intGet <= 30 {
				db, err = sql.Open("mysql", "root:test123@tcp(mysql3:3306)/mysql3")
				if err != nil {
					fmt.Println("Net podklucheniya k DB3", err)
				}
			}

			result, err1 := db.Query("SELECT * FROM `good` WHERE id = " + get + ";")
			if err1 != nil {
				fmt.Println("Нет результата", err)
			}

			goods := []Good{}

			for result.Next() {
				good := Good{}
				result.Scan(&good.Id, &good.Title)
				goods = append(goods, good)
			}

			jsonData, _ := json.Marshal(goods)

			conn.Do("SET", get, string(jsonData))

			cashData = string(jsonData)
		}
		fmt.Fprintf(w, cashData)

	})
	http.ListenAndServe(":8080", nil)
}
