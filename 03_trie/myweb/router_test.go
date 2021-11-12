package myweb

import (
	"fmt"
	"testing"
)
import "github.com/stretchr/testify/assert"

func testNewRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/hello/:name/fd", nil)
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	assert.Equal(t, parsePattern("/p/:name"), []string{"p",":name"})
	assert.Equal(t, parsePattern("/p/name"), []string{"p","name"})
	assert.Equal(t, parsePattern("/p/*"), []string{"p","*"})
	assert.Equal(t, parsePattern("/p/*name/*"), []string{"p", "*name"})
	//log.Println(parsePattern("/p/name/*"))
}

func TestGetRoute(t *testing.T) {
	r := testNewRouter()
	n, ps := r.getRoute("GET", "/hello/myweb")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if ps["name"] != "myweb" {
		t.Fatal("name should be equal to 'myweb'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}

func TestGetRoute2(t *testing.T) {
	r := testNewRouter()
	n1, ps1 := r.getRoute("GET", "/assets/file1.txt")
	ok1 := n1.pattern == "/assets/*filepath" && ps1["filepath"] == "file1.txt"
	if !ok1 {
		t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be file1.txt")
	}

	n2, ps2 := r.getRoute("GET", "/assets/css/test.css")
	ok2 := n2.pattern == "/assets/*filepath" && ps2["filepath"] == "css/test.css"
	if !ok2 {
		t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be css/test.css")
	}
}

func TestGetRoutes(t *testing.T) {
	rt := testNewRouter()
	nodes := rt.getRoutes("GET")
	fmt.Println(nodes)
	for index, value := range nodes {
		fmt.Println(index + 1, value)
	}
	if len(nodes) != 6 {
		t.Fatal("the number of routes should be 6")
	}
}