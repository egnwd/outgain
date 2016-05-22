package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", hello)
	router.POST("/signin", signIn)
	router.StaticFS("/static", http.Dir("static"))
	router.Run()
}

func hello(c *gin.Context) {
	c.HTML(http.StatusOK, "base.tmpl", gin.H{
		"Title":   "Hello",
		"Message": "Hello, World!",
	})
}

func signIn(c *gin.Context) {
	user := fmt.Sprintf("Current User: %s", c.PostForm("idtoken"))
	log.Println(user)
}
