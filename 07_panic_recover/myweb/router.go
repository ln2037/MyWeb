package myweb

import (
	"net/http"
	"strings"
)

type router struct {
	roots 		map[string]*node //存储每种请求的根节点, 如 roots["GET"], roots["POST"]
	handlers 	map[string]HandlerFunc
}

// 解析路由。把/hello/name 解析为字符串切片
func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	sp := strings.Split(pattern, "/")
	for _, item := range sp {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if item[0] == '*' {
			break
		}
	}
	return parts
}

func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//获取所有路由
func (rt *router) getRoutes(method string) []*node {
	root, ok := rt.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.dfs(&nodes)
	return nodes
}

//添加路由，如/hello/:name
func (rt *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	//创建每个方法的根节点
	if _, ok := rt.roots[method]; ok == false {
		rt.roots[method] = &node{}
	}
	//在每个方法的根节点上添加路由
	rt.roots[method].insert(pattern, parts, 0)
	//映射路由和处理函数handler
	rt.handlers[key] = handler
}

//处理路由，返回对应的node,node里的pattern即该路由对应的pattern。返回path值的映射
func (rt *router) getRoute(method string, path string) (*node, map[string]string) {
	pathParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := rt.roots[method]
	if !ok {
		return nil, nil
	}
	nd := root.search(pathParts, 0)
	if nd == nil {
		return nil, nil
	}
	//获取path值的映射
	parts := parsePattern(nd.pattern)
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = pathParts[index]
		} else if part[0] == '*'{
			params[part[1:]] = strings.Join(pathParts[index:], "/")
			break
		}
	}
	return nd, params
}

//处理函数
func (rt *router) handle(context *Context) {
	nd, params := rt.getRoute(context.Method, context.Path)
	if nd != nil {
		context.Params = params
		key := context.Method + "-" + nd.pattern
		context.handlers = append(context.handlers, rt.handlers[key])
		//rt.handlers[key](context)
	} else {
		context.handlers = append(context.handlers, func(context *Context) {
			context.String(http.StatusNotFound,  "404 NOT FOUND: %s\n", context.Path)
		})
		//context.String(http.StatusNotFound,  "404 NOT FOUND: %s\n", context.Path)
	}
	context.Next()
}