package myweb

import "strings"

//trie树的主结构
type node struct {
	//匹配的路由
	pattern string
	//路由中的一段。如/hello/:name中的:name
	part string
	//子节点
	children []*node
	//是否模糊匹配.当当前路由中含:或*的时候为true
	isWild bool
}

//遍历当前路由的子节点，找到pattern不为""的节点
func (nd *node) dfs(nodes *[]*node) {
	if nd.pattern != "" {
		*nodes = append(*nodes, nd)
	}
	for _, child := range nd.children {
		child.dfs(nodes)
	}
}

//插入操作。用于把一个路由插入到trie树中。若不理解可以通过debug看看路由树的结构。实现方式为DFS
//变量不要命名为node。可能会出现意想不到的bug
func (nd *node) insert(pattern string, parts []string, height int) {
	// 递归的终止条件
	if height == len(parts) {
		nd.pattern = pattern
		return
	}
	part := parts[height]
	child := nd.matchChild(part)
	if child == nil {
		child = &node{
			part: part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		nd.children = append(nd.children, child)
	}
	child.insert(pattern, parts, height + 1)
}

// 判断当前节点的子节点能否匹配路由part
func (nd *node) matchChild(part string) *node {
	for _, child := range nd.children {
		if child.part == part || child.isWild == true {
			return child
		}
	}
	return nil
}

// 通过parts匹配对应的pattern, 就是个DFS
func (nd *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(nd.part, "*") {
		// 若pattern为""直接返回nil
		//如/hello/name 匹配/hello/name/class会为nil
		if nd.pattern == "" {
			return nil
		}
		return nd
	}
	part := parts[height]
	//找到当前节点的能够匹配的part的所有子节点
	children := nd.matchChildren(part)

	for _, child := range children {
		if child.part == part || child.isWild == true {
			result := child.search(parts, height + 1)
			//有结构为nil的情况，这种情况不返回
			if result != nil {
				return result
			}
		}
	}
	return nil
}

// 通过part匹配多个child
func (nd *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range nd.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}