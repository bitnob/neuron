package router

import (
	"regexp"
	"strings"
)

type Param struct {
	Key   string
	Value string
}

type node struct {
	path      string
	isParam   bool
	paramName string
	children  []*node
	route     *Route
	regex     *regexp.Regexp
}

func (r *Router) buildRadixTree() *node {
	root := &node{path: "/"}

	for method, routes := range r.routes {
		for _, route := range routes {
			segments := strings.Split(strings.Trim(route.Path, "/"), "/")
			current := root

			for _, segment := range segments {
				if strings.HasPrefix(segment, ":") {
					// Parameter node
					paramName := strings.TrimPrefix(segment, ":")
					child := &node{
						path:      segment,
						isParam:   true,
						paramName: paramName,
					}
					current.children = append(current.children, child)
					current = child
				} else {
					// Static node
					child := &node{path: segment}
					current.children = append(current.children, child)
					current = child
				}
			}

			current.route = &route
		}
	}

	return root
}
