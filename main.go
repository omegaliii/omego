package main

import (
	"net/http"
	"omego"
)

// Engine is the uni handler for all requests
type Engine struct{}

func main() {
	r := omego.New()

	r.GET("/", func(c *omego.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *omego.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *omego.Context) {
		c.JSON(http.StatusOK, omego.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})


	r.Run(":9999")
}