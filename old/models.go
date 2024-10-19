package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

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
	CreatedAt time.Time `json:"created_at"`
}
