package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "myprojectname/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/v1/account/login", login)
	r.POST("api/v1/account/register", register)
	r.GET("/api/v1/account/history", history)
	r.POST("/api/v1/account/chat", chat)
	r.Run(":8080")
}

// login API
// @Summary 使用者登入
// @Description 使用者輸入帳號和密碼進行登入
// @Tags 帳號密碼
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "使用者名稱"
// @Param password formData string true "密碼"
// @Success 200 {object} map[string]interface{} "登入成功"
// @Failure 401 {object} map[string]interface{} "帳號或密碼錯誤"
// @Router /api/v1/account/login [post]
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

// register API
// @Summary 註冊新帳號
// @Description 創建一個新的使用者帳號
// @Tags 帳號密碼
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "使用者名稱"
// @Param password formData string true "密碼"
// @Success 200 {object} map[string]interface{} "創建帳號成功"
// @Failure 500 {object} map[string]interface{} "無法創建帳號"
// @Router /api/v1/account/register [post]
func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	sql := "insert into users (username,password) values ($1,$2)"
	_, err := db.Exec(context.Background(), sql, username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "無法創建帳號密碼",
			"detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "創建帳號成功"})
}

type Chatreq struct {
	Username   string `json:"username"`
	Message    string `json:"message"`
	Gptmessage string
}

// chat API
// @Summary 傳送聊天訊息
// @Description 使用者發送訊息，並接收 GPT 回覆
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body Chatreq true "使用者問題"
// @Success 200 {object} Chatreq "回應gpt 訊息"
// @Failure 400 {object} map[string]interface{} "無效的輸入"
// @Router /api/v1/account/chat [post]
func chat(c *gin.Context) {
	var req Chatreq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的輸入"})
		return
	}
	req.Gptmessage = gptmessage(req.Message)

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

func gptmessage(usermessage string) string {

	apiKey := "AIzaSyDMGjKqIYIl_WTVtM11mNnBe6Z1aPUdtMw" // 將 YOUR_API_KEY 替換為您的 API 金鑰
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	data := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": usermessage,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)

	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)

	}

	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		fmt.Println("No candidates found")
		return ""
	}

	candidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid candidate format")
		return ""
	}

	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid content format")
		return ""
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		fmt.Println("Invalid parts format")
		return ""
	}

	part, ok := parts[0].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid part format")
		return ""
	}

	text, ok := part["text"].(string)
	if !ok {
		fmt.Println("Invalid text format")
		return ""
	}
	return text

}

type ChatHistory struct {
	Username   string `json:"username"`
	Message    string `json:"message"`
	Gptmessage string `json:"gptmessage"`
}

// history API
// @Summary 取得聊天歷史紀錄
// @Description 獲取特定使用者的聊天歷史
// @Tags Chat
// @Produce json
// @Param username query string true "使用者名稱"
// @Success 200 {object} map[string]interface{} "成功回應聊天歷史"
// @Failure 400 {object} map[string]interface{} "請提供使用者名稱"
// @Router /api/v1/account/history [get]
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
