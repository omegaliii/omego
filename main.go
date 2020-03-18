package main

import (
	"net/http"
	"omego"
)

// Engine is the uni handler for all requests
type Engine struct{}

func main() {
	r := omego.New()

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *omego.Context) {
			c.HTML(http.StatusOK, "<h1>Hello omego</h1>")
		})

		v1.GET("/hello", func(c *omego.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *omego.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *omego.Context) {
			c.JSON(http.StatusOK, omego.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")

	r.Run(":9999")
}