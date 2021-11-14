package main

/*
(1)
$ curl -i http://localhost:9999/
HTTP/1.1 200 OK
Date: Mon, 12 Aug 2019 16:52:52 GMT
Content-Length: 18
Content-Type: text/html; charset=utf-8
<h1>Hello myweb</h1>

(2)
$ curl "http://localhost:9999/hello?name=mywebktutu"
hello mywebktutu, you're at /hello

(3)
$ curl "http://localhost:9999/hello/mywebktutu"
hello mywebktutu, you're at /hello/mywebktutu

(4)
$ curl "http://localhost:9999/assets/css/mywebktutu.css"
{"filepath":"css/mywebktutu.css"}

(5)
$ curl "http://localhost:9999/xxx"
404 NOT FOUND: /xxx
*/

import (
	"net/http"

	"MyWeb/03_trie/myweb"
)

func main() {
	r := myweb.New()
	r.GET("/", func(c *myweb.Context) {
		c.HTML(http.StatusOK, "<h1>Hello myweb</h1>")
	})

	r.GET("/hello", func(c *myweb.Context) {
		// expect /hello?name=mywebktutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *myweb.Context) {
		// expect /hello/mywebktutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *myweb.Context) {
		c.JSON(http.StatusOK, myweb.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
