package omego

import (
	"net/http"
)

// Signature for HandlerFunc
type HandlerFunc func(*Context)

// Engine implements the interface of ServeHTTP 
// router [Key: route] [Value: handler function]
type Engine struct {
	router *router
}

// Constructor
func New() *Engine {
	return &Engine{router: newRouter()}
}

// Handle GET request for the path pattern
// @params pattern[string] - path
// @params handler[HandlerFunc] - call back function
//
// @return
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// Handle POST request for the path pattern
// @params pattern[string] - path
// @params handler[HandlerFunc] - call back function
//
// @return
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run the engine in the address
// @params address[string] - server address
func (engine *Engine) Run(address string) (err error) {
	// From net/http:
	// func ListenAndServe(address string, h Handler) error
	// Handler is a interface in net/Http whi
	// Engine must implement ServeHTTP method to be a Handler
	return http.ListenAndServe(address, engine)
}


//################################ 
//######## Helper Methods ######## 
//################################ 

func (engine *Engine) addRoute (method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// To pass an engine to http.ListenAndServe
// Engine must implement ServeHTTP to be a Hanlder(interface from net/http)
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handle(c)
}


