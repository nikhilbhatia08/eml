package codegen

import (
	"os"
	"path/filepath"
	"strings"

	// "fmt"
	"github.com/nikhilbhatia08/eml/parser"
)

// func GenerateCode(root *parser.Node) {
// 	var imports []string
// 	lines := GenerateHtmlCodeFromAST(root)
// }

// THIS IS FRAGILE DO NOT TOUCH
func GenerateHtmlCodeFromAST(node *parser.Node) []string {
	if node == nil {
		return []string{}
	}
	var result []string
	spaces := "\t"
	for i := 0; i < node.TopLeveIndentation; i++ {
		spaces += " "
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
		}else if kv.Key == "styles" {
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
	for _, kv := range node.DirectInfo.Iter() {
		if kv.Key == "path" {
			// fmt.Println("LINKING")
			var linkPath string
			linkPath = kv.Value
			linkPath = spaces + "to={\"" + linkPath + "\"}"
			result = append(result, linkPath)
		} else if kv.Key == "tailwind_styles" {
			// fmt.Println(node.Keyword)
			className := "className=\""
			className += kv.Value
			className += "\""
			result = append(result, spaces + className)
		}
	}
	spaces = spaces[:len(spaces)-1]
	result = append(result, spaces+">")
	spaces += "\t"
	for _, kv := range node.DirectInfo.Iter() {
		// fmt.Println(kv.Key)
		if kv.Key == "content" {
			// fmt.Println(kv.Value)
			content = kv.Value
		}
	} 
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

// This function converts the imports of ehtml to imports of js that are used in react
func GenerateImports(imports []string, routes_length int32) []string {
	// TODO : This function needs proper testing whether the imports are valid or not
	var result []string
	for _, imp := range imports {
		result = append(result, "import " + imp)
	}
	result = append(result, "import { Link } from 'react-router'")
	if routes_length > 0 {
		result = append(result, "import { Route, Routes } from 'react-router'")
	}
	return result
}

// This function is used to generate routes
// This function should validate and transform the routes
func GenerateRoutes(routes []string) ([]string, error) {
	// TODO: This function also needs to check whether the route component is valid or not, should strategize on this well
	if len(routes) == 0 {
		return nil, nil
	}
	var result []string
	result = append(result, "\t<Routes>")
	for _, route := range routes {
		route_tokens := GetRouteTokens(route)
		// if route_tokens[0] != "link" || route_tokens[2] != "to" || len(route_tokens) != 4 {
		// 	return nil, fmt.Errorf("invalid route format: %s", route)
		// }
		// In this function I'm writing magic index numbers for accessing the elements of route_tokens and there should be some strategy
		result = append(result, "\t\t<Route path=\"" + route_tokens[0] + "\" element={<" + GetBasePath(route_tokens[1]) + " />} />")
	}
	result = append(result, "\t</Routes>")
	return result, nil
}

// This is to get the line in the form of tokens
// For example link Comp to /some/path.ehtml
// then the o/p should be ["link", "Comp", "to", "/some/path.ehtml"]
func GetRouteTokens(route string) []string {
	return strings.Split(route, " ")
}

// This function extracts the filetype from a path
func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

// This function extracts only the name of the file
// for example "App" from "components/primary/App.tsx"
func GetBasePath(path string) string {
	return filepath.Base(path)[:len(filepath.Base(path))-len(GetFileExtension(path))]
}

func WriteToFile(path string, lines []string, imports []string, routes []string) error {
	dir := filepath.Dir(path) // "components/primary"

	// Ensure all directories in the path exist
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Write to the file
	fileContents := Combine(path, lines, imports, routes)

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

// This function is used to combine all the ehtml, imports and routes to convert into full react code
func Combine(path string, lines []string, imports []string, routes []string) []string {
	var result []string
	result = append(result, imports...)
	if GetBasePath(path) == "App" {
		result = append(result, "function " + GetBasePath(path) + "() {")
	}else {
		result = append(result, "function " + GetBasePath(path) + "(props) {")
	}
	result = append(result, "\treturn(")
	result = append(result, "\t<>")
	result = append(result, routes...)
	result = append(result, lines...)
	result = append(result, "\t</>")
	result = append(result, "\t)")
	result = append(result, "}")
	result = append(result, "export default " + GetBasePath(path) + ";")
	return result
}