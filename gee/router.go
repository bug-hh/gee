package gee

import (
	"log"
	"net/http"
	"strings"
)

type Router struct {
	roots map[string]*node
	handlers map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

func ParsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) AddRouter(method string, pattern string, handler HandleFunc) {
	log.Printf("Route %s - %s", method, pattern)
	parts := ParsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) GetRoute(method string, path string) (*node, map[string]string) {
	searchParts := ParsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := ParsePattern(n.Pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {
	// 用传入的具体 path 来判断，这个path 属于哪个 Pattern
	// 例如：/hello/bughh 这个 path 就属于 /hello/:name 这个 Pattern
	n, params := r.GetRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		// 所以这个地方用 c.Method + "-" + n.Pattern 来做 key，而不是  c.Method + "-" + c.Path
		key := c.Method + "-" + n.Pattern
		// 这里添加的是用于处理具体请求的 handler，而不是中间件，虽然「处理具体请求的 handler」和 「中间件」都是 HandleFunc 类型
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		log.Printf("not found %s-%s", c.Method, c.Path)
		c.handlers = append(c.handlers, func(ctx *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

