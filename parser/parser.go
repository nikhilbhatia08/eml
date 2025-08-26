package parser

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"github.com/nikhilbhatia08/eml/parser/utils"
)

// We would start parsing from the main function and then
// generate the equivalent code of that

// Recognized tags in Easy HyperText Markup Language (EHTML)
// 1. div
// 2. h1
// Recognized styles in Easy HyperText Markup Language (EHTML)
// 1. text(Represents the size of the text)
// 2. font-bold(Represents bold text)
// 3. underline(Represents underlined text)

func Parser() *Node {
	// Get the lines in the main
	// This can be improved for many files
	lines := ParseFile("main.ehtml")

	root := GenerateAST(lines)
	if root == nil {
		fmt.Println("Compilation failed")
		return nil
	}

	generatedEHTML := root.ConvertToEHTML(root)
	fmt.Println("Generated EHTML:")
	for _, line := range generatedEHTML {
		fmt.Println(line)
	}
	fmt.Println("Compiled successfully")
	return root
}

func ParseFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	// Read the file and parse the contents
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return lines
}

// As soon as the context finishes we need to pop the stack

func GenerateAST(lines []string) *Node {
	// root := &Node{Children: []*Node{}}
	stack := utils.Stack[*Node]{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		spaces := utils.CountSpaces(line)
		if !utils.CheckForCharacter(line) {
			continue
		}
		for !stack.IsEmpty() && stack.Peek().TopLeveIndentation >= spaces {
			// fmt.Println("trace")
			// stack.Pop()
			topNode, err := stack.Pop()
			if err == false {
				// Handle error because stack not being empty could not pop
			}

			if !stack.IsEmpty() {
				// If the stack is not empty, we need to set the parent-child relationship
				if topNode.NodeType == KEYWORD_TYPE {
					parentNode := stack.Peek()
					parentNode.Children = append(parentNode.Children, topNode)
				}else if topNode.NodeType == CONFIG_TYPE {
					parentNode := stack.Peek()
					topNode.Info.Range(func(key string, value *utils.OrderedMap[string, string]) bool {
						parentNode.Info.Set(key, value)
						return true
					})
				}
			}
		}
		lineTokens := GetLineTokens(line)
		if len(lineTokens) == 0 {
			// we need to somehow show error here 
			continue
		}

		if stack.IsEmpty() {
			// If it is empty it means that it has just started
			tokenType := checkToken(lineTokens[0])
			// fmt.Println(lineTokens[0], tokenType)
			if tokenType != KEYWORD_TYPE {
				// We need to show some error here 
			}else {
				newNode := &Node{
					NodeType:     tokenType,
					Keyword:    lineTokens[0],
					Children: []*Node{},
					TopLeveIndentation: spaces,
					Info: utils.NewOrderedMap[string, *utils.OrderedMap[string, string]](),
				}
				stack.Push(newNode)
				// fmt.Println("PUSHING ROOT NODE:", newNode.Keyword, stack.Size(), spaces)
				// root.Children = append(root.Children, newNode)
			}
		}else {
			tokenType := checkToken(lineTokens[0])
			if tokenType == KEYWORD_TYPE {
				// It means that is a keyword
				newNode := &Node{
					NodeType:     tokenType,
					Keyword:    lineTokens[0],
					Children: []*Node{},
					TopLeveIndentation: spaces,
					Info: utils.NewOrderedMap[string, *utils.OrderedMap[string, string]](),
				}
				stack.Push(newNode)
				// fmt.Println("PUSHING ROOT NODE:", newNode.Keyword, stack.Size(), spaces)		
			}else if tokenType == CONFIG_TYPE {
				newNode := &Node{
					NodeType:     tokenType,
					Keyword:    lineTokens[0],
					Children: []*Node{},
					Info: utils.NewOrderedMap[string, *utils.OrderedMap[string, string]](),
					TopLeveIndentation: spaces,
				}
				stack.Push(newNode)
				// fmt.Println("PUSHING CONFIG NODE:", newNode.Keyword, stack.Size(), spaces)
			}else if tokenType == GENERAL_TYPE {
				// fmt.Println("GENERAL TYPE", lineTokens[0])
				if !stack.IsEmpty(){
					topNode := stack.Peek()
					if topNode.NodeType == CONFIG_TYPE {
						if len(lineTokens) < 2 {
							// There should be error handling here
						}else {
							inner, ok := topNode.Info.Get(topNode.Keyword)
							if !ok {
								inner = utils.NewOrderedMap[string, string]() // create new inner ordered map
								topNode.Info.Set(topNode.Keyword, inner)
							}
							var sentence string
							for i:= 1; i < len(lineTokens); i++ {
								if i == len(lineTokens) - 1 {
									sentence += lineTokens[i]
								}else {
									sentence += lineTokens[i] + " "
								}
							}
							// fmt.Println("SETTING:", lineTokens[0], "TO:", sentence, len(lineTokens))
							inner.Set(lineTokens[0], sentence)
						}
					}else {
						// There should be error handling here
					}
				}
			}
		}
	}

	if stack.Size() > 1 {
		for stack.Size() > 1 {
			topNode, err := stack.Pop()
			if err == false {
				// Handle error because stack not being empty could not pop
			}
			if topNode.NodeType == KEYWORD_TYPE {
				// If the top node is a keyword type, we need to set the parent-child relationship
				if !stack.IsEmpty() {
					parentNode := stack.Peek()
					parentNode.Children = append(parentNode.Children, topNode)
				}
			}else if topNode.NodeType == CONFIG_TYPE {
				// If the top node is a config type, we need to merge its info with the parent
				if !stack.IsEmpty() {
					parentNode := stack.Peek()
					topNode.Info.Range(func(key string, value *utils.OrderedMap[string, string]) bool {
						parentNode.Info.Set(key, value)
						return true
					})
				}
			}
		}
	}

	if !stack.IsEmpty() {
		return stack.Peek()
	}

	return stack.Peek()
}

func checkToken(token string) int {
	if Keywords[token] {
		return KEYWORD_TYPE
	}else if config[token] {
		return CONFIG_TYPE
	}
	return GENERAL_TYPE
}

func GetLineTokens(line string) []string {
	var word string
	var lineTokens []string
	for i := 0; i < len(line); i++ {
		if unicode.IsSpace(rune(line[i])) {
			if len(word) > 0 {
				lineTokens = append(lineTokens, word)
				word = ""
			}
			continue
		}else if line[i] == ':' {
			lineTokens = append(lineTokens, word)
			// fmt.Println("PUSHED WORD:", word)
			word = ""
		}else {
			word += string(line[i])
		}
	}
	if len(word) > 0 {
		lineTokens = append(lineTokens, word)
	}
	return lineTokens
}
