package myweb

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

// Engine 实现了ServeHttp接口
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup //存储所有的路由组
	htmlTemplates	*template.Template
	funcMap template.FuncMap //模板函数
}

// 路由组
type RouterGroup struct {
	prefix		string //路由前缀
	middlewares []HandlerFunc //中间件
	engine		*Engine //每个路由组共享一个Engine
	parent		*RouterGroup //路由嵌套
}

//New 创建了Engine的实例，并返回其指针
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default use Logger & Recovery middleware
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// Group 创建一个路由组
// 所有的路由组共享一个 Engine 实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute 添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 添加get请求的路由
func (group *RouterGroup) GET(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("GET", pattern, handlerFunc)
}

// POST 添加POST请求的路由
func (group *RouterGroup) POST(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("POST", pattern, handlerFunc)
}

// 创建静态文件的handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//添加路由前缀
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(context *Context) {
		file := context.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			context.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(context.Writer, context.Request)
	}
}

// Static 添加静态文件的路由
func (group *RouterGroup)Static(relativePath, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	//不必写为urlPattern := path.Join(group.prefix+relativePath, "/*filepath")
	//get方法会自动拼接group.prefix
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

// SetFuncMap 设置模板函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 加载模板
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Run 开启一个http服务
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// HTTP服务
func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		// 若当前url能匹配上对应的路由组，那么添加这个路由组的中间件
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	context := newContext(writer, request)
	//把handlerFunc也添加到中间件里
	context.handlers = middlewares
	context.engine = engine
	engine.router.handle(context)
}