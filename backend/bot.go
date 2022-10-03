package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

type TelegramData struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

type User struct {
	ID        int
	Username  string
	FirstName string
	LastName  string
}

func main() {

	//формируем адрес с данными телеграм
	url := "https://api.telegram.org/bot5527670228:AAGghwOtZWmWfxobTXiggWojqfdDqxCLg3I/getUpdates"

	//создаем подключение к БД
	db, _ := sql.Open("mysql", "root:test123@tcp(mysql:3306)/mysql1")

	//создвем подключение к redis
	conn, _ := redis.Dial("tcp", "redis:6379")

	//имя ключа где лежат кэшированные юзеры
	usersExists := "users"

	//смотрим в Redis есть ли вообще юзеры
	cachedUsers := getRedis(conn, usersExists)

	//смотрим если пустая строка то меняем на пустой json
	if cachedUsers == "" {
		fmt.Println("строка 71 до соединения с бд")
		refreshCahe(usersExists, db, conn)
		fmt.Println("строка 73 после соединения с бд")
	}

	for range time.Tick(time.Second) {

		//отправка GET запроса на адрес
		data, _ := http.Get(url)

		//получаем тело запроса в виде среза байт
		body, _ := ioutil.ReadAll(data.Body)

		//преобразовать ответ к строке чобы проверить визуальтно
		//fmt.Println(string(body))

		//создаем объект типа TelegramData
		telegramData := TelegramData{}

		//распарсить json
		json.Unmarshal(body, &telegramData)

		//достать массив объектов из ответа
		items := telegramData.Result

		updateId := 0
		//в цикле смотрим все обхъекты и достаем данные
		for i := 0; i < len(items); i++ {

			//смотрим в Redis есть ли там в списке такой id юзера
			cachedUsers = getRedis(conn, usersExists)

			//смотрим если пустая строка то меняем на json из Бд
			if cachedUsers == "" {
				_, cachedUsers = refreshCahe(usersExists, db, conn)
			}

			//делаем пустую переменную под юзеров
			usersSlice := []int{}

			//fmt.Println("готовы раскодировать json")

			//заливаем данные из json в переменную usersSlice
			json.Unmarshal([]byte(cachedUsers), &usersSlice)

			//сохдаем юзера
			user := User{ID: items[i].Message.From.ID, Username: items[i].Message.From.Username, FirstName: items[i].Message.From.FirstName, LastName: items[i].Message.From.LastName}
			createUser(usersSlice, user, db, conn, usersExists)

			db.Query("INSERT INTO mysql1.messages(`user_id`,`content`,`c_time`)  VALUES(?,?,?) ", items[i].Message.From.ID, items[i].Message.Text, items[i].Message.Date)

			updateId = items[i].UpdateID
		}

		//увеличиваем updateId на 1 чтобы показывал только новые
		updateId = updateId + 1

		//отправляем updateId
		http.Get(url + "?offset=" + strconv.Itoa(updateId))

	}

	//записать в БД
}

//функция забегает в БД и считывает id юзеров и записывает их в кэш
func refreshCahe(caсheKey string, db *sql.DB, conn redis.Conn) ([]int, string) {

	//зайти в БД и достать всех юзеров
	fmt.Println("До соединения с БД внутри функции")
	result, _ := db.Query("SELECT id FROM `users`")
	fmt.Println("После соединения с БД внутри функции")
	idList := []int{}

	if result != nil {
		for result.Next() {

			id := 0
			result.Scan(&id)
			idList = append(idList, id)
		}
		fmt.Println("после работы с данными")
	}

	//закодировать и записать в кэш
	idJsonData, _ := json.Marshal(idList)
	//conn.Do("SET", caсheKey, string(idJsonData))
	setRedis(conn, caсheKey, string(idJsonData))

	return idList, string(idJsonData)
}

//выясняет есть ли id в списке или нет
func inArrayInt(needle int, haystack []int) bool {

	//перебираем все элементы массива
	for i := 0; i < len(haystack); i++ {

		if haystack[i] == needle {

			return true
		}
	}

	return false
}

//возвращает данные по определенному ключу
func getRedis(conn redis.Conn, key string) string {
	//смотрим в Redis есть ли вообзе юзеры
	data, _ := redis.String(conn.Do("GET", key))
	return data
}

func setRedis(conn redis.Conn, key string, value string) {
	conn.Do("SET", key, value)

}

//проверяет нет ли юзера и его создает
func createUser(list []int, user User, db *sql.DB, conn redis.Conn, key string) bool {
	//если нету то ДЕЛАЕМ ЗАПРОС и ДОПИСЫВАЕМ в REDIS
	if !inArrayInt(user.ID, list) {

		//fmt.Println("такого юзера нет и пишем в базу")

		db.Query("INSERT INTO mysql1.users(`id`,`username`,`first_name`,`last_name`)  VALUES(?,?,?,?) ", user.ID, user.Username, user.FirstName, user.LastName)

		//добавляем id юзера в массив с юзерами
		list = append(list, user.ID)

		//fmt.Println("сделали запрос в базу данных")
		//fmt.Println(usersSlice)

		//кодируем в json юзеров
		jsonData, _ := json.Marshal(list)

		//fmt.Println(string(jsonData))

		//записываем в redis
		setRedis(conn, key, string(jsonData))

		return true

	}

	return false
}
