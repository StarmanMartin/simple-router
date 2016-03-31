package router

import (
	"testing"
    "net/http/httptest"
	"github.com/starmanmartin/simple-router/rtest"
)

func TestMainServer(t *testing.T) {
    resetRouteing()

	router := NewRouter()

	for pathU, handerU := range pathlist {
		router.Get("/tt" + pathU, handerU)
	}
    
    subRouter := NewSubRouter("/martin/starman")
    
    for pathU, handerU := range pathlist {
		subRouter.Get(pathU, handerU)
	}
    
    router = NewXHRRouter()

	for pathU, handerU := range pathlist {
		router.Post("/tt" + pathU, handerU)
	}
    
    subRouter = NewXHRSubRouter("/martin/starman")
    
    for pathU, handerU := range pathlist {
		subRouter.Post(pathU, handerU)
	}
    
    w := httptest.NewRecorder()
    req, _ := rtest.NewGetRequest("http://localhost:8080/tt/hallo", nil)
    router.ServeHTTP(w, req)
    t.Log(w.Code == 200)
    
    w = httptest.NewRecorder()
    req, _ = rtest.NewGetRequest("http://localhost:8080/martin/starman/hallo", nil)
    router.ServeHTTP(w, req)
    t.Log(w.Code == 200)
    
    w = httptest.NewRecorder()
    req, _ = rtest.NewGetRequest("http://localhost:8080/martin/starmans/hallo", nil)
    router.ServeHTTP(w, req)
    t.Log(w.Code == 404)
    
    w = httptest.NewRecorder()
    req, _ = rtest.NewPostReqeust("http://localhost:8080/tt/hallo", nil)
    
    router.ServeHTTP(w, req)
    //t.Log(w.Code)
    
    
    w = httptest.NewRecorder()
    req, _ = rtest.NewXHRPostRequest("http://localhost:8080/tt/hallo", nil)
    router.ServeHTTP(w, req)
    t.Log(w.Code == 200)
    
    w = httptest.NewRecorder()
    req, _ = rtest.NewXHRPostRequest("http://localhost:8080/martin/starman/hallo", nil)
    router.ServeHTTP(w, req)
    t.Log(w.Code == 200)
    
     w = httptest.NewRecorder()
    req, _ = rtest.NewXHRPostRequest("http://localhost:8080/martin/starmans/hallo", nil)
    router.ServeHTTP(w, req)
    t.Log(w.Code == 404)
    
    
}