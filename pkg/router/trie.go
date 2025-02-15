package router

import "strings"

type node struct {
	path     string
	handler  interface{} // Can hold both HandlerFunc and FastHandlerFunc
	children []*node
	isParam  bool
}

func (n *node) find(method, path string) (interface{}, []Param) {
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

	if current.handler != nil {
		return current.handler, params
	}

	return nil, nil
}

func (n *node) insert(method, path string, handler interface{}) {
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
			newNode := &node{
				path:    segment,
				isParam: strings.HasPrefix(segment, ":"),
			}
			current.children = append(current.children, newNode)
			current = newNode
		}
	}
	current.handler = handler
}
