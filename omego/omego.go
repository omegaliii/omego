package omego

import (
	"net/http"
)

// Singature for HandlerFunc
type HandlerFunc func(*Context)

// Engine implements the interface of ServeHTTP 
// router [Key: route] [Value: handler function]
type Engine struct {
	router *Router
}

// Constructor
func New() *Engine {
	return &Engine{router: NewRouter()}
}

func (engine *Engine) addRoute (method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handle(c)
}








