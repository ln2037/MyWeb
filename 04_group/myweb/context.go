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
}

// newContext 返回一个实例的指针
func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     writer,
		Request:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
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

func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	context.Writer.Write([]byte(html))
}