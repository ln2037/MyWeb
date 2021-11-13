package main

import (
	"MyWeb/04_group/myweb"
	"net/http"
)

func main() {
	r := myweb.New()
	r.GET("/index", func(context *myweb.Context) {
		context.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(context *myweb.Context) {
			context.HTML(http.StatusOK, "<h1>hello, myweb</h1>")
		})
		v1.GET("/hello", func(c *myweb.Context) {
			// expect /hello?name=ln2037
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *myweb.Context) {
			// expect /hello/ln2037
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *myweb.Context) {
			c.JSON(http.StatusOK, myweb.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
		v3 := v2.Group("/v3")
		{
			v3.GET("/hello", func(context *myweb.Context) {
				context.HTML(http.StatusOK, "<h1>hello, v2v3</h1>")
			})
			v3.GET("/hello/:name", func(context *myweb.Context) {
				context.HTML(http.StatusOK, "<h1>hello, v2v3</h1>")
				context.String(http.StatusOK, "hello %s, you're at %s\n", context.Param("name"), context.Path)
			})
		}
	}
	r.Run(":9999")
}