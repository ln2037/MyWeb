package main

import (
	"MyWeb/05_middleware/myweb"
	"log"
	"net/http"
	"time"
)

func onlyForV2() myweb.HandlerFunc {
	return func(c *myweb.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func main() {
	r := myweb.New()
	r.Use(myweb.Logger())
	r.GET("/", func(c *myweb.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *myweb.Context) {
			// expect /hello/qwe
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	r.Run(":9999")
}