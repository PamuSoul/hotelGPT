package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

var db *gorm.DB

func main() {

	dsn := "host=127.0.0.1 user=myuser password=mypassword dbname=username port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ 無法連接 postgres:", err)
	}

	db.AutoMigrate(&User{})

	r := gin.Default()

	r.POST("/api/v1/account/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		var user User
		// 查詢資料庫中的使用者
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
			return
		}
		if user.Password == password {
			c.JSON(http.StatusOK, gin.H{"message": "登入成功", "username": username, "token": "aaaaa"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
		}
	})

	r.Run(":8080")
}
