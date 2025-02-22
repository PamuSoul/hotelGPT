package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {

	dsn := "host=127.0.0.1 user=myuser password=mypassword dbname=username port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("❌ 無法連接 PostgreSQL:", err)
	}
	defer db.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/v1/account/login", login)
	r.GET("/api/v1/account/history", history)
	r.POST("/api/v1/account/chat", chat)
	r.Run(":8080")
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var storedPassword string
	sqlStatement := "SELECT password FROM users WHERE username = $1"
	err := db.QueryRow(context.Background(), sqlStatement, username).Scan(&storedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
		return
	}
	if storedPassword == password {
		c.JSON(http.StatusOK, gin.H{"message": "登入成功", "username": username, "token": "AAA"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
	}
}

type Chatreq struct {
	Username   string `json:"username"`
	Message    string `json:"message"`
	Gptmessage string
}

func chat(c *gin.Context) {
	var req Chatreq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的輸入"})
		return
	}
	req.Gptmessage = "奶茶"

	sql := "INSERT INTO history (username, message, gptmessage) VALUES ($1, $2, $3)"
	_, err := db.Exec(context.Background(), sql, req.Username, req.Message, req.Gptmessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "無法儲存訊息",
			"detail": err.Error(), // 顯示 PostgreSQL 回傳的錯誤
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username":   req.Username,
		"message":    req.Message,
		"gptmessage": req.Gptmessage,
	})
}

type ChatHistory struct {
	Username   string `json:"username"`
	Message    string `json:"message"`
	Gptmessage string `json:"gptmessage"`
}

func history(c *gin.Context) {

	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供使用者名稱"})
		return
	}

	sql := "SELECT message, gptmessage FROM history WHERE username = $1 ORDER BY id ASC"
	rows, err := db.Query(context.Background(), sql, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "無法獲取歷史訊息",
			"detail": err.Error()})
		return
	}
	defer rows.Close()

	var histories []ChatHistory
	for rows.Next() {
		var record ChatHistory
		record.Username = username
		if err := rows.Scan(&record.Message, &record.Gptmessage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析歷史訊息"})
			return
		}
		histories = append(histories, record)
	}

	if len(histories) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "無歷史紀錄", "history": []ChatHistory{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": histories})
}
