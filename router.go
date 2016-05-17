package router

import (
	"regexp"
	"strings"
	"github.com/starmanmartin/simple-router/request"
)

const (
	mGet    string = "get"
	mPost   string = "post"
	mPut    string = "put"
	mDelete string = "delete"
	mAll    string = "all"
)

func paramsCopy(src request.Params) request.Params {
	dst := make(request.Params, len(src))

	for k, v := range src {
		dst[k] = v
	}

	return dst
}

var (
	routerList map[string]*routeElement
	splitPath  *regexp.Regexp
)

func init() {
	initTree()
	initUpload()
	resetRouteing()

	splitPath = regexp.MustCompile(`\/\{([^\}]|\}[^\/])+\}\}?`)
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
	list := make(request.Params)
	return findListOfHandlerRec(elem, list, path, xhr)
}

func findListOfHandlerRec(elem *routeElement, params request.Params, path []string, xhr bool) (finalHandler []*finalRouteElement, returnOk bool) {
	if len(path) == 0 {
		return
	}

	if handlerList, tempHandler, ok := elem.getNext(path[0], len(path) == 1); ok {
		finalHandler = make([]*finalRouteElement, 0, len(tempHandler))
		for _, tempElem := range tempHandler {
			if !tempElem.Xhr || xhr {
				tempParams := paramsCopy(params)
				if tempElem.isVariable {
					tempParams[tempElem.route.variableName()] = path[0]
				}

				finalHandler = append(finalHandler, &finalRouteElement{tempElem, &tempParams})
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
	return addElemToTree(handler, routerList[method], prepareRoute(route), xhr)
}

func prepareRoute(route string) []string {
	
	firstSplit := splitPath.Split(route, -1)
	regSplit := splitPath.FindAllString(route, -1)
	if len(firstSplit) > 0 && firstSplit[len(firstSplit)-1] == "" {
		firstSplit = firstSplit[:len(firstSplit)-1]
	}

	var finaleRoute []string
	for i, rp := range firstSplit {
		if rp[:1] == "/" {
			rp = rp[1:]
		}

		finaleRoute = append(finaleRoute, strings.Split(rp, "/")...)
		if i < len(regSplit) {
			if regSplit[i][:1] == "/" {
				regSplit[i] = regSplit[i][1:]
			}

			finaleRoute = append(finaleRoute, regSplit[i])
		}
	}

	if len(finaleRoute) == 0 {
		return []string{""}
	}

	return finaleRoute
}
