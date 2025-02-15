package router

import "strings"

type fastNode struct {
	path     string
	handler  FastHandlerFunc
	children []*fastNode
	isParam  bool
}

func (n *fastNode) find(method, path string) (FastHandlerFunc, []Param) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	params := make([]Param, 0, 8)
	current := n

	for _, segment := range segments {
		found := false
		for _, child := range current.children {
			if child.isParam {
				params = append(params, Param{
					Key:   strings.TrimPrefix(child.path, ":"),
					Value: segment,
				})
				current = child
				found = true
				break
			}
			if child.path == segment {
				current = child
				found = true
				break
			}
		}
		if !found {
			return nil, nil
		}
	}

	return current.handler, params
}

func (n *fastNode) insert(method, path string, handler FastHandlerFunc) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	current := n

	for _, segment := range segments {
		found := false
		for _, child := range current.children {
			if child.path == segment {
				current = child
				found = true
				break
			}
		}
		if !found {
			newNode := &fastNode{
				path:    segment,
				isParam: strings.HasPrefix(segment, ":"),
			}
			current.children = append(current.children, newNode)
			current = newNode
		}
	}
	current.handler = handler
}
