// Package router handels theroutiing
// Author: Maritn Starman
package router

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//HTTPHandler defines Type
type HTTPHandler func(w http.ResponseWriter, r *Request) (bool, error)

var (
	// NotFoundHandler function if path not found
	NotFoundHandler func(w http.ResponseWriter, r *Request)
	// ErrorHandler reponshandler on Error
	ErrorHandler func(err error, w http.ResponseWriter, r *Request)
	// XHRNotFoundHandler function if path not found
	XHRNotFoundHandler func(w http.ResponseWriter, r *Request)
	// XHRErrorHandler reponshandler on Error
	XHRErrorHandler func(err error, w http.ResponseWriter, r *Request)
)

//SubManager manages routs in a sub path
type SubManager struct {
	base string
	xhr  bool
}

//Manager instance of router
type Manager struct {
	*SubManager
}

// Get registers a route for the GET method
func (r *SubManager) Get(route string, handler ...HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mGet, handler, r.xhr)
}

// Post register a route for the POST method
func (r *SubManager) Post(route string, handler ...HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mPost, handler, r.xhr)
}

// Delete register a route for the Delete method
func (r *SubManager) Delete(route string, handler ...HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mDelete, handler, r.xhr)
}

// Put register a route for the Put method
func (r *SubManager) Put(route string, handler ...HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mPut, handler, r.xhr)
}

// All register a route for the All methods
func (r *SubManager) All(route string, handler ...HTTPHandler) {
	route = validateURL(r.base, route)
	if elm, isNew := addNew(route, mAll, handler, r.xhr); isNew {
		addToAll(elm)
	}
}

//Public sets a Fileserver for path
func (r *Manager) Public(path string) string {
	path = validateURL(r.base, path)
	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileServer := http.FileServer(http.Dir(cwd))

	addNew(path+"/*", mGet, []HTTPHandler{func(w http.ResponseWriter, r *Request) (bool, error) {
		fileServer.ServeHTTP(w, r.Request)
		return false, nil
	}}, false)

	return filepath.Join(cwd, path)
}

//UploadPath Registers a upload Path
func (r *Manager) UploadPath(path string, isBuffer bool) HTTPHandler {
	return newUploadPaser(path, isBuffer)
}

func (r *Manager) ServeHTTP(w http.ResponseWriter, httpR *http.Request) {
	req := newRequest(httpR)
	preparedPath := strings.Trim(req.URL.Path, "/")
	method, path := strings.ToLower(req.Method), strings.Split(preparedPath, "/")
	rootElem := routerList[method]
	xhr := httpR.Header.Get("X-Requested-With") == "XMLHttpRequest"
	if finalHandler, ok := findListOfHandler(rootElem, path, xhr); ok {
		for _, tempElem := range finalHandler {
			req.RouteParams = *tempElem.params
			for _, handleFunc := range tempElem.hanlder {
				isNext, err := handleFunc(w, req)

				if err != nil {
                    w.WriteHeader(http.StatusBadRequest)
					if !xhr {
						if ErrorHandler == nil {
							log.Fatal("Error and no Handler")
						} else {
							ErrorHandler(err, w, req)
						}
					} else {
						if XHRErrorHandler == nil {
							log.Fatal("Error and no Handler")
						} else {
							XHRErrorHandler(err, w, req)
						}
					}
				}

				if !isNext {
					return
				}
			}
		}
	}

	if !xhr {
        w.WriteHeader(http.StatusNotFound)
		if NotFoundHandler == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			NotFoundHandler(w, req)
		}
	} else {
		if XHRNotFoundHandler == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			XHRNotFoundHandler(w, req)
		}
	}
}

//NewRouter returns new instance of the Router
func NewRouter() *Manager {
	subRouter := &SubManager{validateURL("/"), false}
	return &Manager{subRouter}
}

//NewSubRouter returns new instance of the Router
func NewSubRouter(root string) *SubManager {
	return &SubManager{validateURL(root), false}
}

//NewXHRRouter returns new instance of the Router
func NewXHRRouter() *Manager {
	subRouter := &SubManager{validateURL("/"), true}
    return &Manager{subRouter}
}

//NewXHRSubRouter returns new instance of the Router
func NewXHRSubRouter(root string) *SubManager {
	return &SubManager{validateURL(root), true}
}

func validateURL(elements ...string) (pathConcated string) {
	for _, elm := range elements {
		if len(elm) > 0 && elm != "/" {
			if elm[:1] != "/" {
				elm = "/" + elm
			}

			if elm[len(elm)-2:] == "/" {
				elm = elm[len(elm)-2:]
			}
		} else {
			elm = ""
		}

		pathConcated += elm
	}

	if len(pathConcated) == 0 {
		pathConcated = "/"
	}

	return
}
