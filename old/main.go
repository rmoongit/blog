package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	connStr := "user=xfx dbname=blog sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Настройка маршрутов
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addPostHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	// Запуск сервера
	log.Println("Сервер запущен на :444")
	log.Fatal(http.ListenAndServe(":444", nil))
}
