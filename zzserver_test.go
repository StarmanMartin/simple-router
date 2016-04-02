package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/starmanmartin/simple-router/request"
	"github.com/starmanmartin/simple-router/rtest"
)

func TestMainServer(t *testing.T) {
	resetRouteing()

	router := NewRouter()

	for pathU, handerU := range pathlist {
		router.Get("/tt"+pathU, handerU)
	}

	subRouter := NewSubRouter("/martin/starman")

	for pathU, handerU := range pathlist {
		subRouter.Get(pathU, handerU)
	}

	router = NewXHRRouter()

	for pathU, handerU := range pathlist {
		router.Post("/tt"+pathU, handerU)
	}

	subRouter = NewXHRSubRouter("/martin/starman")

	for pathU, handerU := range pathlist {
		subRouter.Post(pathU, handerU)
	}

	w := httptest.NewRecorder()
	req, _ := rtest.NewGetRequest("http://localhost:8080/tt/hallo", nil)
	router.ServeHTTP(w, req.Request)
	if w.Code != 200 {
		t.Error("Wrong response code", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = rtest.NewGetRequest("http://localhost:8080/martin/starman/hallo", nil)
	router.ServeHTTP(w, req.Request)
	if w.Code != 200 {
		t.Error("Wrong response code", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = rtest.NewGetRequest("http://localhost:8080/martin/starmans/hallo", nil)
	router.ServeHTTP(w, req.Request)
	if w.Code != 404 {
		t.Error("Wrong response code", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = rtest.NewPostReqeust("http://localhost:8080/tt/hallo", nil)

	router.ServeHTTP(w, req.Request)
	if w.Code != 404 {
		t.Error("Wrong response code", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = rtest.NewXHRPostRequest("http://localhost:8080/tt/hallo", nil)
	router.ServeHTTP(w, req.Request)
	if w.Code != 200 {
		t.Error("Wrong response code", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = rtest.NewXHRPostRequest("http://localhost:8080/martin/starman/hallo", nil)
	router.ServeHTTP(w, req.Request)
	if w.Code != 200 {
		t.Error("Wrong response code", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = rtest.NewXHRPostRequest("http://localhost:8080/martin/starmans/hallo", nil)
	router.ServeHTTP(w, req.Request)
	if w.Code != 404 {
		t.Error("Wrong response code", w.Code)
	}

}

var (
	callList = map[string]int{
		"NotFoundHandler":    0,
		"ErrorHandler":       0,
		"XHRNotFoundHandler": 0,
		"XHRErrorHandler":    0,
	}
)

func Test404Handler(t *testing.T) {

	resetRouteing()

	NotFoundHandler = func(w http.ResponseWriter, r *request.Request) {
		callList["NotFoundHandler"]++
	}

	// XHRNotFoundHandler function if path not found
	XHRNotFoundHandler = func(w http.ResponseWriter, r *request.Request) {
		callList["XHRNotFoundHandler"]++
	}

	router := NewRouter()
	w := httptest.NewRecorder()
	req, _ := rtest.NewGetRequest("http://localhost:8080/martin/starmans/hallo", nil)

	router.ServeHTTP(w, req.Request)

	if !compareErrorCalls(1, 0, 0, 0, callList) {
		t.Error("Not found did not work")
	}

	req, _ = rtest.NewXHRGetRequest("http://localhost:8080/martin/starmans/hallo", nil)

	router.ServeHTTP(w, req.Request)

	if !compareErrorCalls(0, 0, 1, 0, callList) {
		t.Error("XHR not found did not work")
	}
}

func TestNoXhrErrorHandler(t *testing.T) {
	router := NewRouter()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	router.Get("/hallo/error", func(w http.ResponseWriter, r *request.Request) (bool, error) {
		return false, errors.New("Super error")
	})

	req, _ := rtest.NewXHRGetRequest("http://localhost:8080/hallo/error", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req.Request)
}

func TestNoErrorHandler(t *testing.T) {
	router := NewRouter()
	resetRouteing()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	router.Get("/hallo/error", func(w http.ResponseWriter, r *request.Request) (bool, error) {
		return false, errors.New("Super error")
	})

	req, _ := rtest.NewGetRequest("http://localhost:8080/hallo/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req.Request)
}

func TestErrorHandler(t *testing.T) {
	router := NewRouter()
    resetRouteing()
    
    router.Get("/hallo/error", func(w http.ResponseWriter, r *request.Request) (bool, error) {
		return false, errors.New("Super error")
	})

	// ErrorHandler reponshandler on Error
	ErrorHandler = func(err error, w http.ResponseWriter, r *request.Request) {
		callList["ErrorHandler"]++
	}
	// XHRErrorHandler reponshandler on Error
	XHRErrorHandler = func(err error, w http.ResponseWriter, r *request.Request) {
		callList["XHRErrorHandler"]++
	}

	req, _ := rtest.NewGetRequest("http://localhost:8080/hallo/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req.Request)

	if !compareErrorCalls(0, 1, 0, 0, callList) {
		t.Error("Error handler not called")
	}

	req, _ = rtest.NewXHRGetRequest("http://localhost:8080/hallo/error", nil)

	router.ServeHTTP(w, req.Request)

	if !compareErrorCalls(0, 0, 0, 1, callList) {
		t.Error("XHR error handler not called")
	}
}

func compareErrorCalls(NotFoundHandler, ErrorHandler, XHRNotFoundHandler, XHRErrorHandler int, list map[string]int) (retrnVal bool) {
	retrnVal = true

	if NotFoundHandler != list["NotFoundHandler"] {
		retrnVal = false
	}
	if ErrorHandler != list["ErrorHandler"] {
		retrnVal = false
	}
	if XHRErrorHandler != list["XHRErrorHandler"] {
		retrnVal = false
	}
	if XHRNotFoundHandler != list["XHRNotFoundHandler"] {
		retrnVal = false
	}

	list["XHRNotFoundHandler"] = 0
	list["XHRErrorHandler"] = 0
	list["ErrorHandler"] = 0
	list["NotFoundHandler"] = 0

	return
}
