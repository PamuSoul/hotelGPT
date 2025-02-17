package main

import (
	"context"
	"log"
	"net/http"

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

	r.POST("/api/v1/account/login", func(c *gin.Context) {
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
	})

	r.Run(":8080")
}
