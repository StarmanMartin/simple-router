package router

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/starmanmartin/simple-router/request"
)

var (
	elementIndex                    = 0
	regexpRegexp, regexpValueRegexp *regexp.Regexp
)

func initTree() {
	regexpRegexp = regexp.MustCompile("^\\{.+\\}$")
	regexpValueRegexp = regexp.MustCompile("^\\{([^\\}\\=]+\\=)?(.+)\\}$")
}

type routePath string

func (rp routePath) match(toMatch string) bool {
	return rp.isEqualToString("*") || rp.routeName() == toMatch
}

func (rp routePath) isWildcard() bool {
	idx := strings.Index(string(rp), "*")
	return idx >= 0 && idx == len(rp)-1
}

func (rp routePath) routeName() string {
	if rp.isVariable() {
		return rp.variableName()
	}

	if rp.isWildcard() {
		return string(rp)[:len(rp)-1]
	}

	return string(rp)
}

func (rp routePath) isVariable() bool {
	return strings.Index(string(rp), ":") == 0
}

func (rp routePath) isRegexp() bool {
	return regexpRegexp.MatchString(string(rp))
}

func (rp routePath) getRegexpKeys() (string, bool, string) {
	keys := regexpValueRegexp.FindAllStringSubmatch(string(rp), -1)
	var isVaribale bool
	if len(keys[0][1]) > 0 {
		keys[0][1] = keys[0][1][:len(keys[0][1])-1]
		isVaribale = true
	}

	return keys[0][2], isVaribale, keys[0][1]
}

func (rp routePath) variableName() string {
	return strings.Split(string(rp), ":")[1]
}

func (rp routePath) isEqualToString(text string) bool {
	return string(rp) == text
}

type routeElement struct {
	route                                    routePath
	routeAsRegx                              *regexp.Regexp
	next                                     []*routeElement
	hanlder                                  []HTTPHandler
	isVariable, isFinal, isMatchAll, isRegex bool
	Xhr                                      bool
	index                                    int
}

type finalRouteElement struct {
	*routeElement
	params *request.Params
}

func newRouteElement(routeName string, xhr bool) (tempElement *routeElement) {
	route := routePath(routeName)
	elementIndex++
	tempElement = &routeElement{route, nil, make([]*routeElement, 0), nil, false, false, false, false, xhr, elementIndex}
	if route.isVariable() {
		tempElement.isVariable = true
	}

	if route.isWildcard() {
		tempElement.isMatchAll = true
	} else if route.isRegexp() {
		tempElement.isRegex = true
		routeRegexp, hasVariable, varName := route.getRegexpKeys()
		tempElement.routeAsRegx = regexp.MustCompile(routeRegexp)

		if hasVariable {
			tempElement.isVariable = true
		}
		
		tempElement.route = routePath(":" + varName + ":" + string(tempElement.route))
	}

	return
}

func (b *routeElement) getNext(pathElem string, isLast bool) ([]*routeElement, []*routeElement, bool) {
				
	var nextList, finalIndex = make([]*routeElement, 0, len(b.next)), 0
	
	for _, p := range b.next {
		if p.matchRegexp(pathElem) || p.route.match(pathElem) || (p.isVariable && !p.isRegex) {
			if p.isFinal && isLast || p.isMatchAll {
				nextList = append([]*routeElement{p}, nextList...)
				finalIndex++
			} else {
				nextList = append(nextList, p)
			}
		}
	}

	return nextList[finalIndex:], nextList[:finalIndex], len(nextList) != 0
}

func (b *routeElement) matchRegexp(path string) (ok bool){
	if !b.isRegex {
		return
	}

	return b.routeAsRegx.MatchString(path)
}

func (b *routeElement) isEqualToPath(path string) (ok bool){
	if b.isRegex {
		comp := strings.Split(string(b.route), ":")[2]
		return comp == path
	}

	return b.route.isEqualToString(path)
}

func (b *routeElement) String() string {
	childs := ""
	for _, el := range b.next {
		childs = fmt.Sprintf("%s, %s", childs, el)
	}

	if len(childs) == 0 {
		return fmt.Sprintf("%s-%t", b.route, b.isFinal)
	}

	return fmt.Sprintf("%s-%t (%s)", b.route, b.isFinal, childs)
}

func addToAll(elm *routeElement) {
	for idx, root := range routerList {
		if idx != mAll {
			root.next = append(root.next, elm)
		}
	}
}

func addElemToTree(handler []HTTPHandler, treeNode *routeElement, routeParts []string, xhr bool) (*routeElement, bool) {
	if len(routeParts) == 0 {
		if treeNode.isFinal {
			panic(fmt.Sprintf("Double route!! Last Element"))
		}
		treeNode.hanlder, treeNode.isFinal = handler, true
		return treeNode, false
	}

	for _, nextNode := range treeNode.next {
		if nextNode.isEqualToPath(routeParts[0]) {
			addElemToTree(handler, nextNode, routeParts[1:], xhr)
			return nextNode, false
		}
	}

	newNode := newRouteElement(routeParts[0], xhr)
	treeNode.next = append(treeNode.next, newNode)
	addElemToTree(handler, treeNode, routeParts, xhr)
	return newNode, true
}
