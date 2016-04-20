package router

import (
	"testing"
)

func TestRouter(t *testing.T) {
	router := NewRouter()
	if message, isErr := testManager(router.SubManager, "/", false); isErr {
		t.Error(message)
	}
}

func TestSubRouter(t *testing.T) {
	router := NewSubRouter("test")
	if message, isErr := testManager(router, "/test", false); isErr {
		t.Error(message)
	}
}

func TestXHRRouter(t *testing.T) {
	router := NewXHRRouter()
	if message, isErr := testManager(router.SubManager, "/", true); isErr {
		t.Error(message)
	}
}

func TestSubXHRRouter(t *testing.T) {
	router := NewXHRSubRouter("test")
	if message, isErr := testManager(router, "/test", true); isErr {
		t.Error(message)
	}
}

func TestValidRoute(t *testing.T) {
	route := validateURL("//test")
	if route != "/test" {
		t.Error("validateURL(\"//test\") -> response wrong: ", route)
	}
	route = validateURL("/test/")
	if route != "/test" {
		t.Error("validateURL(\"/test/\") -> response wrong: ", route)
	}
	route = validateURL("//test//")
	if route != "/test" {
		t.Error("validateURL(\"//test//\") -> response wrong: ", route)
	}

	route = validateURL("")
	if route != "/" {
		t.Error("validateURL(\"/\") -> response wrong: ", route)
	}

	route = validateURL("test")
	if route != "/test" {
		t.Error("validateURL(\"test\") -> response wrong: ", route)
	}

	route = validateURL("test1", "", "/test2", "test3/", "/test4/", "a")
	if route != "/test1/test2/test3/test4/a" {
		t.Error("validateURL(\"test1\", \"\", \"/test2\", \"test3/\", \"/test4/\", \"a\") -> response wrong: ", route)
	}
}

//--------------------------------------------------------------------------------------

func testManager(m *SubManager, base string, isXhr bool) (message string, isErr bool) {
	if m.base != base {
		return "base wrong:" + m.base, true
	}

	if m.xhr != isXhr {
		return "XHR error", true
	}

	return "", false
}