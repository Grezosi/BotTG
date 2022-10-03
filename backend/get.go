package main

import (
	
	"fmt"
	"database/sql"
  _ "github.com/go-sql-driver/mysql"
  	"net/http"
	"encoding/json"
)

type Item struct {
UserId int
Username string
Content string
CreationTime int

}

type Info struct {
	Count int
	Messages []MessageT
}

type MessageT struct {
	Content string
	Date int
}


func main() {

	//подключаемся к Бд и достаем юзеров
	db, _ := sql.Open("mysql", "root:test123@tcp(localhost:3310)/mysql1")

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {

		//зайти в БД и получить инфо о всех юзерах и сообщениях
		result, _ := db.Query("SELECT m.`user_id`, m.`content`, m.`c_time`, u. `username` FROM `messages` m LEFT JOIN `users` u ON u.id=m.user_id")
		items := []Item{}
		for result.Next(){
			item := Item{}
			result.Scan(&item.UserId,&item.Content,&item.CreationTime,&item.Username)
			items = append(items,item)
		}

		//пересобираем данные в нужный нам формат
		list := map[int]Info{}
		for i :=0; i < len(items); i++{

			//смотрим есть ли элемент с таким ключом в таблицу
			_, exist := list[items[i].UserId]
			//если есть, увеличиваем значение на 1
			if exist {
				//смотрим инфо по данному юзеру
				info := list[items[i].UserId]
				//в инфо смотрим сколько у юзера сообщ и увеличиваем на 1
				info.Count = info.Count +1 
				//в инфо смотрим н амассив сообщ и дописываем туда новое сообщ
				info.Messages = append(info.Messages, MessageT{Content: items[i].Content, Date: items[i].CreationTime})
				//обновляем информацию по данному юзеру в списке
				list[items[i].UserId] = info
			} else {
				//если нет, создаем и пишем значение 1
				//messages := []string{items[i].Content}
				messages := []MessageT{}
				messages = append(messages, MessageT{Content: items[i].Content, Date: items[i].CreationTime})
				info := Info{Count: 1, Messages: messages}
				list[items[i].UserId] = info
			}
			
		}

		jsonData, _ := json.Marshal(list)

		fmt.Fprintf(w, string(jsonData))

	})

	http.ListenAndServe(":8901", nil)
}