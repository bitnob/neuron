package middleware

import (
	"fmt"
	"runtime"
)

func Recover() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					fmt.Printf("panic: %v\n\n%s", r, buf[:n])
					c.JSON(500, Error{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: "Internal server error",
					})
				}
			}()
			return next(c)
		}
	}
}
