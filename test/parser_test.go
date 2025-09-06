package test

import (
	"fmt"
	"testing"

	"github.com/nikhilbhatia08/eml/parser"
)

func equal(a, b []string) bool {
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

func TestLines(t *testing.T) {
	fmt.Println("Testing Lines")
	line := "	h1:"
	tokens := parser.GetLineTokens(line)
	expected := []string{"h1"}
	if !equal(tokens, expected) {
		t.Errorf("Expected %v, but got %v", expected, tokens)
	}

	line = "		styles:"
	tokens = parser.GetLineTokens(line)
	expected = []string{"styles"}
	if !equal(tokens, expected) {
		t.Errorf("Expected %v, but got %v", expected, tokens)
	}

	line = "		text: 3xl"
	tokens = parser.GetLineTokens(line)
	expected = []string{"text", "3xl"}
	if !equal(tokens, expected) {
		t.Errorf("Expected %v, but got %v", expected, tokens)
	}
	// Passing Message
}

func TestParser(t *testing.T) {
	fmt.Println("Testing Parser")
	lines := []string{
		"import:",
		"	Comp from ../hehe",
		"	Comp3 from ./fellow",
		"			",
		"router:",
		"	/Comp3: Comp3",
		"div:",
		"	h1:",
		"		tailwind_styles: text-3xl font-bold underline",
		"							",
		"		content: Hello World",
	}

	root, imports, routes := parser.GenerateAST(lines)
	if root == nil {
		t.Fatal("Failed to parse EHTML")
	}

	// ehtml := root.ConvertToEHTML(root)
	// expected := []string{
	// 	"div:",
	// 	"\th1:",
	// 	"\t\ttailwind_styles: text-3xl font-bold underline",
	// 	"\t\tcontent: Hello World",
	// }

	expected_imports := []string {
		"Comp from \"../hehe\"",
		"Comp3 from \"./fellow\"",
	}

	expected_routes := []string{
		// CONFUSING
		"/Comp3 Comp3", // This will not have a colon because it is enriched
	}
	// if !equal(ehtml, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, ehtml)
	// 	for i := 0; i < len(ehtml); i++ {
	// 		fmt.Println(ehtml[i])
	// 	}
	// 	for i := 0; i < len(expected); i++ {
	// 		fmt.Println(expected[i])
	// 	}
	// }

	if !equal(imports, expected_imports) {
		t.Errorf("Expected %v, but got %v", expected_imports, imports)
		for i := 0; i < len(imports); i++ {
			fmt.Println(imports[i])
		}
		for i := 0; i < len(expected_imports); i++ {
			fmt.Println(expected_imports[i])
		}
	}
	if !equal(routes, expected_routes) {
		t.Errorf("Expected %v, but got %v", expected_routes, routes)
		for i := 0; i < len(routes); i++ {
			fmt.Println(routes[i])
		}
		for i := 0; i < len(expected_routes); i++ {
			fmt.Println(expected_routes[i])
		}
	}
}

// need to add more test cases to line tokens for more coverage
func TestLineTokens(t *testing.T) {
	line := "\t/Second: Second"
	tokens := parser.GetLineTokens(line)
	expected := []string{"/Second", "Second"}
	if !equal(tokens, expected) {
		t.Errorf("Expected %v, but got %v", expected, tokens)
	}
}

// We need to add more tests for generating ast for more scenarios
