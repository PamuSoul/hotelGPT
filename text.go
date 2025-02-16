package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var dbusername string = "ABC"
var dbpassword string = "123"

func main() {
	r := gin.Default()

	r.POST("/api/v1/account/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == dbusername && password == dbpassword {
			c.JSON(http.StatusOK, gin.H{"message": "登入成功", "username": username, "token": "aaaaa"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
		}
	})

	r.Run(":8080")
}
