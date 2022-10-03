package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {

	//подключаемся к Бд и достаем юзеров
	db, _ := sql.Open("mysql", "root:test123@tcp(localhost:3310)/mysql1")

	usersId := getIds(db)

	//создаем канал для записи ответа
	ch := make(chan string, len(usersId))

	//в цикле смотрим их id и отправляем сообщение(в отдельном потоке каждое)
	for i := 0; i < len(usersId); i++ {

		go sendMessage(usersId[i], "hello", ch)
	}

	all := []string{}
	//считываем данные из канала
	for i := 0; i < len(usersId); i++ {
		//fmt.Println(<-ch)
		all = append(all, <-ch)
	}

	fmt.Println(all)

	//пишем в файл
	file, _ := os.Create("send-log.txt")
	defer file.Close()

	jsonStr, _ := json.Marshal(all)
	file.WriteString(string(jsonStr))
}

func sendMessage(chatId int, text string, ch chan string) {

	resp, _ := http.Get("https://api.telegram.org/bot5527670228:AAGghwOtZWmWfxobTXiggWojqfdDqxCLg3I/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=" + text)

	//прочитать тело запроса
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- string(body)
}

func getIds(db *sql.DB) []int {
	result, _ := db.Query("SELECT `id` FROM `users`")
	fmt.Println("вывожу результат")
	usersId := []int{}
	for result.Next() {
		id := 0
		result.Scan(&id)
		usersId = append(usersId, id)
	}

	return usersId

}
