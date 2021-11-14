package main

import (
	"MyWeb/07_panic_recover/myweb"
	"net/http"
)

func main() {
	r := myweb.Default()
	r.GET("/", func(context *myweb.Context) {
		context.String(http.StatusOK, "Hello, ln2037!\n")
	})
	r.GET("/panic", func(context *myweb.Context) {
		names := []string{"asdf"}
		context.String(http.StatusOK, names[1])
	})
	r.Run(":9999")
}
