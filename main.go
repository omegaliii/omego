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
		c.HTML(http.StatusOK, "<h1>Hello omego</h1>")
	})

	r.GET("/hello", func(c *omego.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *omego.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *omego.Context) {
		c.JSON(http.StatusOK, omego.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}