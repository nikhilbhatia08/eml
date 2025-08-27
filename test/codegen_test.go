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

func TestCodeGenTest(t *testing.T) {
	lines := []string{
		"div:",
		"	h1:",
		"		styles:",
		"			text: 3xl",
		"			font: bold",
		"			underline: true",
		"							",
		"		content:",
		"			v: Hello World",
	}

	root := parser.GenerateAST(lines)
	if root == nil {
		t.Fatal("Failed to parse EHTML")
	}

	html := codegen.GenerateHtmlCodeFromAST(root)
	expected := []string{
		"<div",
		">",
		"	<h1",
		"		className=\"text-3xl font-bold underline\"",
		"	>",
		"		Hello World",
		"	</h1>",
		"</div>",
	}
	if !equal1(html, expected) {
		t.Errorf("Expected %v, but got %v", expected, html)
	}
	fmt.Println("Code Generation Test Passed")
}
