package main

import (
	"github.com/gocraft/web"
)

// Context is initalised for each user by the middleware and passed
// to the URL handlers where it can be inspected/updated.
type Context struct {
}

// NewRouter creates a web router for the main site
func NewRouter() *web.Router {

	r := web.New(Context{})

	// Setup the middleware
	r.Middleware(web.StaticMiddleware("public"))
	r.Middleware((*Context).RequestLogger)

	return r

}
