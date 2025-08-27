package codegen

import (
	"os"
	"path/filepath"
	// "fmt"
	"github.com/nikhilbhatia08/eml/parser"
)

// func GenerateCode(root *parser.Node) {
// 	var imports []string
// 	lines := GenerateHtmlCodeFromAST(root)
// }


// THIS IS FRAGILE DO NOT TOUCH
func GenerateHtmlCodeFromAST(node *parser.Node) []string {
	// At first there should be imports so we don't suppport it for now
	if node == nil {
		return []string{}
	}
	var result []string
	spaces := "\t"
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

func GenerateImports(imports []string) []string {
	var result []string
	for _, imp := range imports {
		result = append(result, "import " + imp)
	}
	return result
}

func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

func GetBasePath(path string) string {
	return filepath.Base(path)[:len(filepath.Base(path))-len(GetFileExtension(path))]
}

func WriteToFile(path string, lines []string, imports []string) error {
	dir := filepath.Dir(path) // "components/primary"

	// Ensure all directories in the path exist
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Write to the file
	fileContents := Combine(path, lines, imports)

	output := ""
	for _, line := range fileContents {
		output += line + "\n"
	}
	// for _, line := range fileContents {
	err = os.WriteFile(path, []byte(output), 0644)
	if err != nil {
		panic(err)
	}
	// }
	// fmt.Println("File written to:", path)
	return nil
}

func Combine(path string, lines []string, imports []string) []string {
	var result []string
	result = append(result, imports...)
	result = append(result, "function " + GetBasePath(path) + "(props) {")
	result = append(result, "\treturn(")
	result = append(result, lines...)
	result = append(result, "\t)")
	result = append(result, "}")
	result = append(result, "export default " + GetBasePath(path) + ";")
	return result
}