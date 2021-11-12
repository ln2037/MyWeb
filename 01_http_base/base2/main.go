package main

import (
	"MyWeb/01_http_base/base2/myweb"
	"fmt"
	"net/http"
)

func main() {
	router := myweb.New()
	router.GET("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
	})
	router.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	router.Run(":9999")
}

