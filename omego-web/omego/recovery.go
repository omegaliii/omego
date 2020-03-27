package omego

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	// skip first 3 caller: callers, trace, defer func
	n := runtime.Callers(3, pcs[:]) 

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		// get the function
		fn := runtime.FuncForPC(pc)
		// get file name and line
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}