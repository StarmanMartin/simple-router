// Package router handels theroutiing
// Author: Maritn Starman
package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//HTTPHandler defines Type
type HTTPHandler func(w http.ResponseWriter, r *Request)

//SubManager manages routs in a sub path
type SubManager struct {
	base string
}

//Manager instance of router
type Manager struct {
	*SubManager
}

// Get registers a route for the GET method
func (r *SubManager) Get(route string, handler HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mGet, handler)
}

// Post register a route for the POST method
func (r *SubManager) Post(route string, handler HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mPost, handler)
}

// Delete register a route for the Delete method
func (r *SubManager) Delete(route string, handler HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mDelete, handler)
}

// Put register a route for the Put method
func (r *SubManager) Put(route string, handler HTTPHandler) {
	route = validateURL(r.base, route)
	addNew(route, mPut, handler)
}

// All register a route for the All methods
func (r *SubManager) All(route string, handler HTTPHandler) {
	route = validateURL(r.base, route)
	if elm, isNew := addNew(route, mAll, handler); isNew {
		addToAll(elm)
	}
}

//Public sets a Fileserver for path
func (r *Manager) Public(path string) string {
	path = validateURL(r.base, path)
	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileServer := http.FileServer(http.Dir(cwd))

	addNew(path+"/*", mGet, func(w http.ResponseWriter, r *Request) {
		fileServer.ServeHTTP(w, r.Request)
	})

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
	if finalHandler, ok := findListOfHandler(rootElem, path); ok {
		for _, tempElem := range finalHandler {
			req.RouteParams = *tempElem.params
			tempElem.hanlder(w, req)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

//NewRouter returns new instance of the Router
func NewRouter() *Manager {
	subRouter := &SubManager{validateURL("/")}
	return &Manager{subRouter}
}

//NewSubRouter returns new instance of the Router
func NewSubRouter(root string) *SubManager {
	return &SubManager{validateURL(root)}
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
