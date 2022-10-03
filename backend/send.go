package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	//подключаемся к Бд и достаем юзеров
	db, _ := sql.Open("mysql", "root:test123@tcp(localhost:3310)/mysql1")

	http.HandleFunc("/api/sendMessage", func(w http.ResponseWriter, r *http.Request) {

		message := r.URL.Query().Get("text")

		usersId := getIds(db)

		//в цикле смотрим их id и отправляем сообщение(в отдельном потоке каждое)
		for i := 0; i < len(usersId); i++ {

			go sendMessage(usersId[i], message)
		}

	})

	http.HandleFunc("/api/sendPhoto", func(w http.ResponseWriter, r *http.Request) {

		message := r.URL.Query().Get("photo")

		usersId := getIds(db)

		//в цикле смотрим их id и отправляем сообщение(в отдельном потоке каждое)
		for i := 0; i < len(usersId); i++ {

			go sendPhoto(usersId[i], message)
		}

	})

	http.HandleFunc("/api/sendLocation", func(w http.ResponseWriter, r *http.Request) {

		latitude, _ := strconv.ParseFloat(r.URL.Query().Get("latitude"), 32)
		longitude, _ := strconv.ParseFloat(r.URL.Query().Get("longitude"), 32)

		usersId := getIds(db)

		//в цикле смотрим их id и отправляем сообщение(в отдельном потоке каждое)
		for i := 0; i < len(usersId); i++ {

			go sendLocation(usersId[i], float32(latitude), float32(longitude))
		}

	})

	http.ListenAndServe(":8900", nil)

}

func sendMessage(chatId int, text string) {
	http.Get("https://api.telegram.org/bot5527670228:AAGghwOtZWmWfxobTXiggWojqfdDqxCLg3I/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=" + text)
}

func sendPhoto(chatId int, photo string) {

	http.Get("https://api.telegram.org/bot5527670228:AAGghwOtZWmWfxobTXiggWojqfdDqxCLg3I/sendPhoto?chat_id=" + strconv.Itoa(chatId) + "&photo=" + photo)
}

func sendLocation(chatId int, latitude float32, longitude float32) {

	http.Get("https://api.telegram.org/bot5527670228:AAGghwOtZWmWfxobTXiggWojqfdDqxCLg3I/sendLocation?chat_id=" + strconv.Itoa(chatId) + "&latitude=" + fmt.Sprintf("%f", latitude) + "&longitude=" + fmt.Sprintf("%f", longitude))
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
