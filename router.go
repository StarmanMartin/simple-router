package router

import (
	"strings"
)

const (
	mGet    string = "get"
	mPost   string = "post"
	mPut    string = "put"
	mDelete string = "delete"
	mAll    string = "all"
)

type params map[string]string

func paramsCopy(src params) params {
	dst := make(params, len(src))

	for k, v := range src {
		dst[k] = v
	}

	return dst
}

var routerList map[string]*routeElement

func init() {
	resetRouteing()
}

func resetRouteing() {
	routerList = make(map[string]*routeElement)
	routerList[mGet] = newRouteElement(mGet, false)
	routerList[mPost] = newRouteElement(mPost, false)
	routerList[mPut] = newRouteElement(mPut, false)
	routerList[mDelete] = newRouteElement(mDelete, false)
	routerList[mAll] = newRouteElement(mAll, false)
}

func findListOfHandler(elem *routeElement, path []string, xhr bool) ([]*finalRouteElement, bool) {
	list := make(params)
	return findListOfHandlerRec(elem, list, path, xhr)
}

func findListOfHandlerRec(elem *routeElement, params params, path []string, xhr bool) (finalHandler []*finalRouteElement, returnOk bool) {
	if len(path) == 0 {
		return
	}

	if handlerList, tempHandler, ok := elem.getNext(path[0], len(path) == 1); ok {
		finalHandler = make([]*finalRouteElement, len(tempHandler))
		for idx, tempElem := range tempHandler {
			if !tempElem.Xhr || xhr {
				tempParams := paramsCopy(params)
				if tempElem.isVariable {
					tempParams[tempElem.route.variableName()] = path[0]
				}

				finalHandler[idx] = &finalRouteElement{tempElem, &tempParams}
			}
		}

		for _, tempElem := range handlerList {
			tempParams := paramsCopy(params)
			if tempElem.isVariable {
				tempParams[tempElem.route.variableName()] = path[0]
			}

			if tempHandlerNext, okNext := findListOfHandlerRec(tempElem, tempParams, path[1:], xhr); okNext {
				finalHandler = append(finalHandler, tempHandlerNext...)
			}
		}
	}

	returnOk = len(finalHandler) != 0
	sortRoutingList(finalHandler)
	return
}

func addNew(route string, method string, handler []HTTPHandler, xhr bool) (*routeElement, bool) {
	if route[:1] == "/" {
		route = route[1:]
	}

	return addElemToTree(handler, routerList[method], strings.Split(route, "/"), xhr)
}
