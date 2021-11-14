package myweb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request
	//url的路径
	Path string
	//请求方法
	Method string
	//返回的状态码
	StatusCode int
	//参数
	Params map[string]string
	//中间件
	handlers []HandlerFunc
	index int
	engine *Engine
}

// newContext 返回一个实例的指针
func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     writer,
		Request:    req,
		Path:       req.URL.Path,
		Method:     req.Method,
		index: 		-1,
	}
}

// Next 调用下一个handlerFunc
func (context *Context) Next() {
	//index代表当前执行到了第几个中间件
	context.index++
	size := len(context.handlers)
	for ; context.index < size; context.index++ {
		context.handlers[context.index](context)
	}
}

func (context *Context) Param(key string) string {
	value, _ := context.Params[key]
	return value
}

// PostForm 获得表单里的元素
func (context *Context) PostForm(key string) string {
	return context.Request.FormValue(key)
}

// Query 获得url里的参数
func (context *Context) Query(key string) string {
	return context.Request.URL.Query().Get(key)
}

// Status 设置HTTP状态码
func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

func (context *Context) Fail(code int, err string) {
	context.index = len(context.handlers)
	context.JSON(code, H{"message": err})
}

// SetHeader 设置Header
func (context *Context) SetHeader(key, value string) {
	context.Writer.Header().Set(key, value)
}

func (context *Context) 	String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (context *Context) JSON(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(code)
	//获取编码器
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), 500)
	}
}

func (context *Context) Data(code int, data []byte) {
	context.Status(code)
	context.Writer.Write(data)
}

// HTML template render
func (context *Context) HTML(code int, name string, data interface{}) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	if err := context.engine.htmlTemplates.ExecuteTemplate(context.Writer, name, data); err != nil {
		context.Fail(500, err.Error())
	}
}