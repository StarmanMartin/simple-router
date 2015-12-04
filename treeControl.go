package router

import (
	"fmt"
	"strings"
)

var elementIndex = 0

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

func (rp routePath) variableName() string {
	if strings.Index(string(rp), ":") == 0 {
		return string(rp)[1:]
	}

	return string(rp)
}

func (rp routePath) isEqualToString(text string) bool {
	return string(rp) == text
}

func (rp routePath) isEqualToPath(text routePath) bool {
	return string(rp) == string(text)
}

type routeElement struct {
	route      routePath
	next       []*routeElement
	hanlder    HTTPHandler
	isVariable bool
	isFinal    bool
	isMatchAll bool
	index      int
}

type finalRouteElement struct {
	*routeElement
	params *params
}

func newRouteElement(routeName string) (tempElement *routeElement) {
	route := routePath(routeName)
	elementIndex++
	tempElement = &routeElement{route, make([]*routeElement, 0), nil, false, false, false, elementIndex}
	if route.isVariable() {
		tempElement.isVariable = true
	}

	if route.isWildcard() {
		tempElement.isMatchAll = true
	}

	return
}

func (b routeElement) getNext(pathElem string, isLast bool) ([]*routeElement, []*routeElement, bool) {
	var nextList, finalIndex = make([]*routeElement, 0, len(b.next)), 0
	for _, p := range b.next {
		if p.route.match(pathElem) || p.isVariable {
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

func (b routeElement) String() string {
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

func addElemToTree(handler HTTPHandler, treeNode *routeElement, routeParts []string) (*routeElement, bool) {
	if len(routeParts) == 0 {
		if treeNode.isFinal {
			panic(fmt.Sprintf("Double route!! Last Element"))
		}
		treeNode.hanlder, treeNode.isFinal = handler, true
		return treeNode, false
	}

	for _, nextNode := range treeNode.next {
		if nextNode.route.isEqualToString(routeParts[0]) {
			addElemToTree(handler, nextNode, routeParts[1:])
			return nextNode, false
		}
	}

	newNode := newRouteElement(routeParts[0])
	treeNode.next = append(treeNode.next, newNode)
	addElemToTree(handler, treeNode, routeParts)
	return newNode, true
}
