package myweb

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (router *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) handle(context *Context) {
	key := context.Method + "-" + context.Path
	if handler, ok := router.handlers[key]; ok {
		handler(context)
	} else {
		context.String(http.StatusNotFound,  "404 NOT FOUND: %s\n", context.Path)
	}
}