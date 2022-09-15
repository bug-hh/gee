package main

import (
	"fmt"
	"github.com/Gee/gee"
	"log"
	"reflect"
)

func newTestRounter() *gee.Router {
	r := gee.NewRouter()
	r.AddRouter("GET", "/", nil)
	r.AddRouter("GET", "/hello/:name", nil)
	r.AddRouter("GET", "/hello/b/c", nil)
	r.AddRouter("GET", "/hi/:name", nil)
	r.AddRouter("GET", "/assets/*filepath", nil)

	return r
}

func TestParsePattern() {
	ok := reflect.DeepEqual(gee.ParsePattern("/p/:name"), []string{"p", ":name"})
	if !ok {
		log.Fatalf("test parsePattern failed")
	}
}

func TestGetRoute() {
	r := newTestRounter()
	n, ps := r.GetRoute("GET", "/hello/bughh")

	if n == nil {
		log.Fatalf("nil shouldn't be returned")
	}

	if n.Pattern != "/hello/:name" {
		log.Fatalf("should match /hello/:name")
	}

	if ps["name"] != "bughh" {
		log.Fatalf("name should be equal to bughh")
	}

	fmt.Printf("match path: %s, params['name']: %s\n", n.Pattern, ps["name"])
}

func main() {
	TestGetRoute()
}


