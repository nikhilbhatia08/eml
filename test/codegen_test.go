package test

import (
	"fmt"
	"testing"

	"github.com/nikhilbhatia08/eml/codegen"
	"github.com/nikhilbhatia08/eml/parser"
)

func equal1(a, b []string) bool {
	if len(a) != len(b) {
		fmt.Println("len unequal")
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			fmt.Printf("At line %d: Expected %q, but got %q\n", i+1, b[i], a[i])
			return false
		}
	}
	return true
}

func TestExtension(t *testing.T) {
	path := "components/primary/button.tsx"
	expected := ".tsx"
	if got := codegen.GetFileExtension(path); got != expected {
		t.Errorf("GetBasePath(%q) = %q; want %q", path, got, expected)
	}
}

func TestBasePath(t *testing.T) {
	path := "components/primary/button.tsx"
	expected := "button"
	if got := codegen.GetBasePath(path); got != expected {
		t.Errorf("GetBasePath(%q) = %q; want %q", path, got, expected)
	}
}

func TestGenerateRoutes(t *testing.T) {
	routes := []string{
		"/Second Second",
	}
	expected := []string{
		"\t<Routes>",
		"\t\t<Route path=\"/Second\" element={<Second />} />",
		"\t</Routes>",
	}
	if got, _ := codegen.GenerateRoutes(routes); !equal1(got, expected) {
		t.Errorf("GenerateRoutes(%v) = %v; want %v", routes, got, expected)
	}
}

func TestCodeGenTest(t *testing.T) {
	lines := []string{
		"div:",
		"	h1:",
		"		tailwind_styles: text-3xl font-bold underline",
		"							",
		"		content: Hello World",
	}

	root, _, _ := parser.GenerateAST(lines)
	if root == nil {
		t.Fatal("Failed to parse EHTML")
	}

	html := codegen.GenerateHtmlCodeFromAST(root)
	expected := []string{
		"\t<div",
		"\t>",
		"\t	<h1",
		"\t		className=\"text-3xl font-bold underline\"",
		"\t	>",
		"\t		Hello World",
		"\t	</h1>",
		"\t</div>",
	}
	if !equal1(html, expected) {
		t.Errorf("Expected %v, but got %v", expected, html)
	}
	fmt.Println("Code Generation Test Passed")
}

func TestCombine(t *testing.T) {
	lines := []string {
		"\t<div",
		"\t>",
		"\t	<h1",
		"\t		className=\"text-3xl font-bold underline\"",
		"\t	>",
		"\t		Hello World",
		"\t	</h1>",
		"\t</div>",
	}

	imports := []string {
		"import Comp1 from \"./fellow\"",
		"import { Route, Routes } from 'react-router'",
	}

	routes := []string {
		"\t<Routes>",
		"\t\t<Route path=\"/Second\" element={<Second />} />",
		"\t</Routes>",
	}

	path := "/some/Path.js"

	combinedOutput := codegen.Combine(path, lines, imports, routes)

	expectedCombinedOutput := []string {
		"import Comp1 from \"./fellow\"",
		"import { Route, Routes } from 'react-router'",
		"function " + codegen.GetBasePath(path) + "(props) {",
		"\treturn(",
		"\t<>",
		"\t<Routes>",
		"\t\t<Route path=\"/Second\" element={<Second />} />",
		"\t</Routes>",
		"\t<div",
		"\t>",
		"\t	<h1",
		"\t		className=\"text-3xl font-bold underline\"",
		"\t	>",
		"\t		Hello World",
		"\t	</h1>",
		"\t</div>",
		"\t</>",
		"\t)",
		"}",
		"export default " + codegen.GetBasePath(path) + ";",
	}

	if !equal1(combinedOutput, expectedCombinedOutput) {
		t.Errorf("Expected %v, but got %v", expectedCombinedOutput, combinedOutput)
	}

}
