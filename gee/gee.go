package gee

import (
	// go template 的用法 https://www.cnblogs.com/f-ck-need-u/p/10053124.html
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
	// 用来加载模板
	htmlTemplates *template.Template
	// 模板渲染函数
	funcMap template.FuncMap
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

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 这里一定加 *，用指针
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
	
}
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		log.Printf("file: %s", file)
		if _, err := fs.Open(file); err != nil {
			log.Printf("not found")
			ctx.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	log.Printf("group.prefix: %s", group.prefix)
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
}

// 把这里写成了 group RouterGroup 而不是写指针类型，导致中间件不起作用，添加失败
func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
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
	var middlewares []HandleFunc
	// 当我们接收到一个具体请求时，要判断该请求适用于哪些中间件
	// 这里只遍历了 group 的 middleware，还没有遍历全局的 middleware，也就是说，全局的 middleware 还没有加入 context 中
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	// 得到中间件列表后，赋值给 c.handlers
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
