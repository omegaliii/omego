package omego

import (
	"strings"
	"log"
	"net/http"
	"path"
	"html/template"
)

// Signature for HandlerFunc
type HandlerFunc func(*Context)

// Engine implements the interface of ServeHTTP
// router [Key: route] [Value: handler function]
type Engine struct {
	*RouterGroup //Embeddding, all the methods of RouterGroup are available on Engine
	router       *router
	groups       []*RouterGroup

	//HTML render
	htmlTemplates *template.Template 
	funcMap       template.FuncMap
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a same Engine instance
}

// Constructor of engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Serve static files
// @params relativePath[string] - relative path to the static file
// @parmas root[string] - root path
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET handlers
	group.GET(urlPattern, handler)
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

//######################################
//########### Add route  ###############
//######################################

// Create a new group
// @params prefix[string] - prefix path
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	// All groups share the same Engine instance
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix, // nesting
		parent: group,                 // nesting
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Append middle to the group.middlewares
// @params middlewares[...HandlerFunc] - a list of handler function
func (group *RouterGroup) Use(midlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, midlewares...)
}

// Handle POST request for the path pattern via Group
// @params pattern[string] - path
// @params handler[HandlerFunc] - call back function
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// Handle POST request for the path pattern via Group
// @params pattern[string] - path
// @params handler[HandlerFunc] - call back function
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//#####################################
//######## HTML render Methods ########
//#####################################

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//################################
//######## Helper Methods ########
//################################

// Add route via engine
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Register Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

// Add route via group
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Register Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// To pass an engine to http.ListenAndServe
// Engine must implement ServeHTTP to be a Handler(interface from net/http)
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// When a request comes in, append middlewares to context from the groups which prefix in the URL path
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// Create static handler function
// @params relativePath[string] - relative path to the file
// @params fs[http.FileSystem]  
//
// @return [callback function] 
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// Map absolutePath to http.FileSystem
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}