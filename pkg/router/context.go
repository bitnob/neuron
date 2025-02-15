package router

import (
	"encoding/json"
	"net/http"
)

func (c *Context) Status(code int) *Context {
	c.Response.WriteHeader(code)
	return c
}

func (c *Context) JSON(code int, data interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(code)
	return json.NewEncoder(c.Response).Encode(data)
}

func (c *Context) Blob(code int, contentType string, data []byte) error {
	c.Response.Header().Set("Content-Type", contentType)
	c.Response.WriteHeader(code)
	_, err := c.Response.Write(data)
	return err
}

func (c *Context) NoContent(code int) error {
	c.Response.WriteHeader(code)
	return nil
}

// NewContext creates a new Context instance
func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Request:  r,
		Response: w,
		store:    make(map[string]interface{}),
	}
}

// String sends a string response
func (c *Context) String(code int, s string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(code)
	_, err := c.Response.Write([]byte(s))
	return err
}
