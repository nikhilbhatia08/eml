package parser

import (
	"github.com/nikhilbhatia08/eml/parser/utils"
)

type Node struct {
	Keyword  string
	Info     *utils.OrderedMap[string, *utils.OrderedMap[string, string]] // Stores metadata about the node
	NodeType int
	TopLeveIndentation int
	Children []*Node
}

func (n* Node) ConvertToEHTML(node *Node) []string {
	if node == nil {
		return []string{}
	}
	var result []string
	spaces := ""
	for i := 0; i < node.TopLeveIndentation; i++ {
		spaces += "\t"
	}
	result = append(result, spaces + node.Keyword + ":")
	// Convert the node to EHTML format
	spaces += "\t"
	for _, kv := range node.Info.Iter() {
		result = append(result, spaces+kv.Key+":")
		spaces += "\t"
		for _, innerKV := range kv.Value.Iter() {
			result = append(result, spaces+innerKV.Key+": "+innerKV.Value)
		}
		spaces = spaces[:len(spaces)-1] // Reduce indentation
	}
	for _, child := range node.Children {
		result = append(result, child.ConvertToEHTML(child)...)
	}
	return result
}
