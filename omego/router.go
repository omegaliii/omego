package omego

import (
    "net/http"
    "strings"
)

// roots    [Key: method] [Value: path node]
// handlers [Key: route]  [Value: handler function]
type router struct {
    roots map[string]*node
    handlers map[string]HandlerFunc
}

// Constructor
// @return *router
func newRouter() *router {
    return &router{
        roots:    make(map[string]*node),
        handlers: make(map[string]HandlerFunc),
    }
}

// Parse the pattern
// @param pattern[string] - path
//
// @return []string - parsed string array of path
func parsePattern(pattern string) []string {
	splittedPattern := strings.Split(pattern, "/")

    parts := make([]string, 0)

    // Iterative through splittedPattern as well as validate each item
	for _, item := range splittedPattern {
		if item != "" {
            parts = append(parts, item)
            
            // Only one * is allowed
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// Add route to the router
// @param method[string] - http verb
// @param pattern[string] - path
// @param handler[HandlerFunc] - handler function for the method-path
//
// @return
func (router *router) addRoute (method string, pattern string, handler HandlerFunc) {
    parts := parsePattern(pattern)

    key := method + "-" + pattern

    _, ok := router.roots[method]

    // Create a new node if it is not exist
    if !ok {
        router.roots[method] = &node{}
    }

    router.roots[method].insert(pattern, parts, 0)
    router.handlers[key] = handler
}

// Get route node from the router as well collect the parameters
// @param method[string] - http verb
// @param pattern[string] - path
//
// @return node - leaf node of the path
// @return map[string]string - parameters of the path
func (router *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := router.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
            // Handle ":"
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
            }
            // Handle "*"
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break // Only allow one *
			}
		}
		return n, params
	}

	return nil, nil
}

// Handle the context
// @param c[Context] 
func (router *router) handle(c *Context) {
    node, params := router.getRoute(c.Method, c.Path)

	// Append the user defined handler function into c.handlers to handle later
    if node != nil {
		c.Params = params
        key := c.Method + "-" + node.pattern
		c.handlers = append(c.handlers, router.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	c.Next()
}








