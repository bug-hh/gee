package main

import (
	"github.com/Gee/gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandleFunc {
	return func(ctx *gee.Context) {
		t := time.Now()
		ctx.Fail(500, "Internal server error")
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}

func main() {
	engine := gee.New()
	engine.Use(gee.Logger()) // global middleware
	engine.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v2 := engine.Group("/v2")
	v2.Use(onlyForV2())  // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	engine.Run(":9999")
}
