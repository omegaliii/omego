package omego

import (
	"net/http"
)

// handlers [Key: route] [Value: handler function]
type Router struct {
	handlers map[string]HandlerFunc
}

// Constructor
func NewRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

func (router *Router) addRoute (method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *Router) handle(c *Context) {
	key := c.Method + "-" + c.Path

	if handler, ok := router.handlers[key]; ok{
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}








