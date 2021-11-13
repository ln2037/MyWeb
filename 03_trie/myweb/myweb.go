package myweb

import (
	"net/http"
)

type HandlerFunc func(*Context)

// Engine 实现了ServeHttp接口
type Engine struct {
	router *router
}



//New 创建了Engine的实例，并返回其指针
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

// addRoute 添加路由
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	engine.addRoute("GET", pattern, handlerFunc)
}

func (engine *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	engine.addRoute("POST", pattern, handlerFunc)
}

// Run 开启一个http服务
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := newContext(writer, request)
	engine.router.handle(context)
}