package codegen

import (
	"github.com/nikhilbhatia08/eml/parser"
)

// func GenerateCode(root *parser.Node) {
// 	lines := GenerateHtmlCodeFromAST(root)
// }

func GenerateHtmlCodeFromAST(node *parser.Node) []string {
	// At first there should be imports so we don't suppport it for now
	if node == nil {
		return []string{}
	}
	var result []string
	spaces := ""
	for i := 0; i < node.TopLeveIndentation; i++ {
		spaces += "\t"
	}
	result = append(result, spaces + "<" + node.Keyword)
	// Convert the node to EHTML format
	content := ""
	spaces += "\t"
	for _, kv := range node.Info.Iter() {
		if kv.Key == "content" {
			for _, innerKv := range kv.Value.Iter() {
				content = innerKv.Value
			}
		}else {
			if kv.Key == "styles" {
				className := "className=\""
				for _, innerKv := range kv.Value.Iter() {
					if innerKv.Value == "true" {
						className += innerKv.Key + " "
					}else {
						className += innerKv.Key + "-" + innerKv.Value + " "
					}
				}
				className = className[:len(className)-1] + "\""
				result = append(result, spaces+className)
			}
		}
	}
	spaces = spaces[:len(spaces)-1]
	result = append(result, spaces+">")
	spaces += "\t"
	if len(content) > 0 {
		result = append(result, spaces + content)
	}
	for _, child := range node.Children {
		result = append(result, GenerateHtmlCodeFromAST(child)...)
	}
	spaces = spaces[:len(spaces)-1]
	result = append(result, spaces + "</" + node.Keyword + ">")
	return result
}