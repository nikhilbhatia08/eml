package parser

import (
	"fmt"
	"testing"
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
	tokens := GetLineTokens(line)
	expected := []string{"h1"}
	if !equal(tokens, expected) {
		t.Errorf("Expected %v, but got %v", expected, tokens)
	}

	line = "		styles:"
	tokens = GetLineTokens(line)
	expected = []string{"styles"}
	if !equal(tokens, expected) {
		t.Errorf("Expected %v, but got %v", expected, tokens)
	}

	line = "		text: 3xl"
	tokens = GetLineTokens(line)
	expected = []string{"text", "3xl"}
	if !equal(tokens, expected) {
		t.Errorf("Expected %v, but got %v", expected, tokens)
	}
	// Passing Message
}

func TestParser(t *testing.T) {
	fmt.Println("Testing Parser")
	lines := []string{
		"div:",
		"	h1:",
		"		styles:",
		"			text: 3xl",
		"			font-bold: true",
		"			underline: false",
		"							",
		"		content:",
		"			v: Hello World",
	}

	root := GenerateAST(lines)
	if root == nil {
		t.Fatal("Failed to parse EHTML")
	}

	ehtml := root.ConvertToEHTML(root)
	expected := []string{
		"div:",
		"\th1:",
		"\t\tstyles:",
		"\t\t\ttext: 3xl",
		"\t\t\tfont-bold: true",
		"\t\t\tunderline: false",
		"\t\tcontent:",
		"\t\t\tv: Hello World",
	}
	if !equal(ehtml, expected) {
		t.Errorf("Expected %v, but got %v", expected, ehtml)
		for i := 0; i < len(ehtml); i++ {
			fmt.Println(ehtml[i])
		}
		for i := 0; i < len(expected); i++ {
			fmt.Println(expected[i])
		}
	}
	// Passing message
}

// We need to add more tests for generating ast for more scenarios
