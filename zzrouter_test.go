package router

import (
	"net/http"
	"strings"
	"testing"
    "net/http/httptest"
	"github.com/starmanmartin/simple-router/rtest"
)

var (
	resivedList [6]int
	pathlist    = map[string]HTTPHandler{
		"/*": func(w http.ResponseWriter, r *Request) (bool, error) {
			resivedList[0]++
			return false, nil
		},
		"/": func(w http.ResponseWriter, r *Request) (bool, error) {
			resivedList[1]++
			return false, nil
		},
		"/:martin/du": func(w http.ResponseWriter, r *Request) (bool, error) {
			resivedList[2]++
			return false, nil
		},
		"/hallo*": func(w http.ResponseWriter, r *Request) (bool, error) {
			resivedList[3]++
			return false, nil
		},
		"/hallo/*": func(w http.ResponseWriter, r *Request) (bool, error) {
			resivedList[4]++
			return false, nil
		},
		"/hallo/:name": func(w http.ResponseWriter, r *Request) (bool, error) {
			resivedList[5]++
			return false, nil
		},
	}

	orders = map[string][6]int{
		"":         {1, 1, 0, 0, 0, 0},
		"hallo":    {1, 0, 0, 1, 0, 0},
		"hallo/du": {1, 0, 1, 1, 1, 1},
	}
)

func TestDoubleRoute(t *testing.T) {
    resetRouteing()
	router := NewRouter()
    
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("The code did not panic")
        }
    }()

    // The following is the code under test
    router.Get("/", nil)
    router.Get("/", nil)
}

func TestGetRouting(t *testing.T) {
	resetRouteing()

	router := NewRouter()

	for pathU, handerU := range pathlist {
		router.Get(pathU, handerU)
	}

	routingTestUtil(t, router.SubManager, mGet)
}

func TestPostRouting(t *testing.T) {
	resetRouteing()
	router := NewRouter()

	for pathU, handerU := range pathlist {
		router.Post(pathU, handerU)
	}

	routingTestUtil(t, router.SubManager, mPost)

}

func TestDelRouting(t *testing.T) {
	resetRouteing()
	router := NewRouter()

	for pathU, handerU := range pathlist {
		router.Delete(pathU, handerU)
	}

	routingTestUtil(t, router.SubManager, mDelete)

}

func TestPutRouting(t *testing.T) {
	resetRouteing()
	router := NewRouter()

	for pathU, handerU := range pathlist {
		router.Put(pathU, handerU)
	}

	routingTestUtil(t, router.SubManager, mPut)
}

func TestAllRouting(t *testing.T) {
	resetRouteing()
	router := NewRouter()

	for pathU, handerU := range pathlist {
		router.All(pathU, handerU)
	}

	routingTestUtil(t, router.SubManager, mGet)
	routingTestUtil(t, router.SubManager, mPost)
	routingTestUtil(t, router.SubManager, mDelete)
	routingTestUtil(t, router.SubManager, mPut)
}

func TestRedirect(t *testing.T) {
    rh, _ := rtest.NewFileUploadRequest("http://localhost:80", make(map[string]string), "image", "test/image.png")   
    r := newRequest(rh)
    res := httptest.NewRecorder()
    r.Redirect(res, "/test")
    if "/test" != res.Header().Get("Location") {
        t.Error("Redirect did not work")
    }
}

//--------------------------------------------------------------------------------------

func routingTestUtil(t *testing.T, router *SubManager, method string) {
	for pathU, list := range orders {
		resivedList = [...]int{0, 0, 0, 0, 0, 0}
		if finalHandler, ok := findListOfHandler(routerList[method], strings.Split(pathU, "/"), false); ok {
			for _, finalHanderTemp := range finalHandler {
                finalHanderTemp.routeElement.String()
				up := *finalHanderTemp.params
				if valP, isOK := up["martin"]; isOK && valP != "hallo" {
					t.Error("Param wrong", valP)
				}

				for _, handerTemp := range finalHanderTemp.hanlder {
					handerTemp(nil, nil)
				}
			}

		}

		for listIndex := range resivedList {
			if resivedList[listIndex] != list[listIndex] {
				t.Error("wrong routing:", pathU, listIndex)
			}
		}
	}
}
