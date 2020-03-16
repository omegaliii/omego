package omego

import (
	"net/http"
)

// handlers [Key: route] [Value: handler function]
type router struct {
	handlers map[string]HandlerFunc
}

// Constructor
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (router *router) addRoute (method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path

	if handler, ok := router.handlers[key]; ok{
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}








