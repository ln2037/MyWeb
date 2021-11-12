package main

import (
	"MyWeb/02_context/myweb"
	"net/http"
)

func main() {
	router := myweb.New()
	router.GET("/", func(context *myweb.Context) {
		context.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	router.GET("/hello", func(context *myweb.Context) {
		context.String(http.StatusOK, "hello %s, you're at %s\n", context.Query("name"), context.Path)
	})

	router.Run(":9999")
}

