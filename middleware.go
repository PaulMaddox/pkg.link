package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/gocraft/web"
)

// ContentTypeJSON sets the content header for JSON responses
func ContentTypeJSON(res web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	res.Header().Set("Content-Type", "application/json")
	next(res, req)
}

// RequestLogger prints a pretty JSON representation of incoming requests
func (c *Context) RequestLogger(res web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	log.Printf("***************************************************************************************")
	log.Printf("* %s: %s", req.Method, req.URL.Path)
	log.Printf("***************************************************************************************")
	for header := range req.Header {
		log.Printf("* %s: %s", header, req.Header[header])
	}

	content := req.Header.Get("Content-Type")
	if content == "text/plain" || content == "application/json" || content == "application/json; charset=utf-8" {

		body, _ := ioutil.ReadAll(req.Body)
		if len(body) > 0 {
			log.Printf("* %s", string(body[:]))
		}

		// As we've now drained req.Body, we need
		// to refill it so other middleware/handlers
		// don't get an empty body.
		restore := bytes.NewReader(body)
		req.Body = ioutil.NopCloser(restore)

	}

	log.Printf("***************************************************************************************")

	next(res, req)
}
