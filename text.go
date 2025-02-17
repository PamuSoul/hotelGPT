package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	dsn := "host=127.0.0.1 user=myuser password=mypassword dbname=username port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ 無法連接 PostgreSQL:", err)
	}

	r := gin.Default()

	r.POST("/api/v1/account/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var storedPassword string
		sqlStatement := "SELECT password FROM users WHERE username = $1"
		err := db.QueryRow(sqlStatement, username).Scan(&storedPassword)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
			return
		}
		if storedPassword == password {
			c.JSON(http.StatusOK, gin.H{"message": "登入成功", "username": username, "token": "AAA"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
		}
	})

	r.Run(":8080")
}
