package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

var db *sql.DB
var jwtKey = []byte("j4K/gU37OlJqkh26Z/kkLTv6Dy9x36vuaiBFXUBM/Vo") // Замените на ваш секретный ключ

// Структура для хранения данных пользователя
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Структура для хранения токена
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Структура для поста
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// Обработчик для главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Обработчик для добавления поста
func addPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Проверка JWT
	tokenString := r.Header.Get("Authorization")
	claims := &Claims{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	username := claims.Username

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post.Username = username
	// Вставка поста в базу данных
	_, err = db.Exec("INSERT INTO posts (title, content, username) VALUES ($1, $2, $3)", post.Title, post.Content, post.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения постов
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, content, username, created_at FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Обработчик для регистрации пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Сохранение пользователя в базе данных
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Обработчик для входа пользователя
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка пользователя в базе данных
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", user.Username).Scan(&storedPassword)
	if err != nil || storedPassword != user.Password {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	// Создание токена
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка токена пользователю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func main() {
	var err error
	connStr := "user=xfx password=passxfx dbname=blog sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Закрытие соединения с базой данных при завершении работы

	// Настройка маршрутов
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addPostHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/posts", getPostsHandler)

	// Запуск сервера
	log.Println("Сервер запущен на :444")
	log.Fatal(http.ListenAndServe(":444", nil))
}
