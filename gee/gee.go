package gee

import (
	"log"
	"net/http"
)

type HandleFunc func(ctx *Context)

type RouterGroup struct {
	prefix string
	middlewares []HandleFunc
	parent *RouterGroup
	engine *Engine
}

type Engine struct {
	// 这是一种嵌套类型，类似 Java/Python 等语言的继承。这样 Engine 就可以拥有 RouterGroup 的属性了。
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{
		router: NewRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		parent:      group,
		engine:      group.engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group RouterGroup) addRoute(method string, comp string, handle HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.addRoute(method, pattern, handle)
}

func (group RouterGroup) GET(pattern string, handle HandleFunc) {
	group.addRoute("GET", pattern, handle)
}

func (group RouterGroup) POST(pattern string, handle HandleFunc) {
	group.addRoute("POST", pattern, handle)
}

func (engine *Engine) addRoute(method string, pattern string, handler HandleFunc) {
	engine.router.AddRouter(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandleFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandleFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
