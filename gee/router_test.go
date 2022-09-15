package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRounter() *Router {
	r := NewRouter()
	r.AddRouter("GET", "/", nil)
	r.AddRouter("GET", "/hello/:name", nil)
	r.AddRouter("GET", "/hello/b/c", nil)
	r.AddRouter("GET", "/hi/:name", nil)
	r.AddRouter("GET", "/assets/*filepath", nil)

	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(ParsePattern("/p/:name"), []string{"p", ":name"})
	if !ok {
		t.Fatal("test ParsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRounter()
	n, ps := r.GetRoute("GET", "/hello/bughh")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.Pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "bughh" {
		t.Fatal("name should be equal to bughh")
	}

	fmt.Printf("match path: %s, params['name']: %s\n", n.Pattern, ps["name"])
}




