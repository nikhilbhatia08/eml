package parser

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"unicode"

	// "fmt"

	"github.com/nikhilbhatia08/eml/parser/utils"
)

// We would start parsing from the main function and then
// generate the equivalent code of that

// Recognized tags in Easy HyperText Markup Language (EHTML)
// 1. div
// 2. h1
// 3. Link for routing purpose
// Recognized styles in Easy HyperText Markup Language (EHTML)
// 1. text(Represents the size of the text)
// 2. font-bold(Represents bold text)
// 3. underline(Represents underlined text)

func Parser() []Target {
	// Get the lines in the main
	// This can be improved for many files
	var targets []Target

	rootDir := "./" // starting path

	// WalkDir will recursively go through all subdirectories
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only parse files that end with .ehtml
		if !d.IsDir() && filepath.Ext(d.Name()) == ".ehtml" {
			// fmt.Println("Parsing:", path)

			lines := ParseFile(path) // parse file lines

			root, imports, routes := GenerateAST(lines)
			if root == nil {
				// There should be proper error message for this one 
				// fmt.Println("Compilation failed for:", path)
				return nil // skip, don’t stop parsing others
			}

			// fmt.Println("Compiled successfully:", path)

			targets = append(targets, Target{
				Path: path[:len(path)-6] + ".js",
				Root: root,
				Imports: imports,
				Routes: routes,
			})
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return targets
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

func GenerateAST(lines []string) (*Node, []string, []string) {
	// root := &Node{Children: []*Node{}}
	stack := utils.Stack[*Node]{}
	var imports []string
	// we need to make sure that the routes array is only of the App.ehtml file
	var routes []string
	for _, line := range lines {
		if line == "" {
			continue
		}
		spaces := utils.CountSpaces(line)
		// fmt.Println(spaces)
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
			}else if topNode.NodeType == IMPORT_TYPE {
				imports = topNode.Imports
			} else if topNode.NodeType == ROUTER_TYPE {
				routes = topNode.Router
			}
		}
		lineTokens := GetLineTokens(line)
		if len(lineTokens) == 0 {
			// we need to somehow show error here 
			continue
		}

		// fmt.Println(lineTokens)
		if stack.IsEmpty() {
			// If it is empty it means that it has just started
			tokenType := checkToken(lineTokens[0])
			// fmt.Println(lineTokens[0], tokenType)
			if tokenType == IMPORT_TYPE {
				newNode := &Node{
					NodeType:     tokenType,
					Keyword:    lineTokens[0],
					Children: []*Node{},
					TopLeveIndentation: spaces,
					Info: utils.NewOrderedMap[string, *utils.OrderedMap[string, string]](),
					DirectInfo: utils.NewOrderedMap[string, string](),
					Imports: make([]string, 0),
				}
				stack.Push(newNode)
			}else if tokenType == ROUTER_TYPE {
				newNode := &Node {
					NodeType:     tokenType,
					Keyword:    lineTokens[0],
					Children: []*Node{},
					TopLeveIndentation: spaces,
					Info: utils.NewOrderedMap[string, *utils.OrderedMap[string, string]](),
					DirectInfo: utils.NewOrderedMap[string, string](),
					Router: make([]string, 0),
				}
				stack.Push(newNode)
			}else if tokenType != KEYWORD_TYPE {
				// We need to show some error here 
			}else {
				newNode := &Node{
					NodeType:     tokenType,
					Keyword:    lineTokens[0],
					Children: []*Node{},
					TopLeveIndentation: spaces,
					Info: utils.NewOrderedMap[string, *utils.OrderedMap[string, string]](),
					DirectInfo: utils.NewOrderedMap[string, string](),
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
					DirectInfo: utils.NewOrderedMap[string, string](),
				}
				// fmt.Println(lineTokens[0])
				stack.Push(newNode)
				// fmt.Println("PUSHING ROOT NODE:", newNode.Keyword, stack.Size(), spaces)		
			}else if tokenType == CONFIG_TYPE {
				newNode := &Node{
					NodeType:     tokenType,
					Keyword:    lineTokens[0],
					Children: []*Node{},
					DirectInfo: utils.NewOrderedMap[string, string](),
					Info: utils.NewOrderedMap[string, *utils.OrderedMap[string, string]](),
					TopLeveIndentation: spaces,
				}
				stack.Push(newNode)
				// fmt.Println("PUSHING CONFIG NODE:", newNode.Keyword, stack.Size(), spaces)
			}else if tokenType == GENERAL_TYPE {
				// fmt.Println("GENERAL TYPE", lineTokens[0])
				if !stack.IsEmpty(){
					topNode := stack.Peek()
					if topNode.NodeType == KEYWORD_TYPE {
						var content_string string
						for i := 1; i < len(lineTokens); i++ {
							if i == len(lineTokens) - 1 {
								content_string += lineTokens[i];
							}else {
								content_string += lineTokens[i] + " ";
							}
						}
						topNode.DirectInfo.Set(lineTokens[0], content_string)
					}else if topNode.NodeType == CONFIG_TYPE {
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
					}else if topNode.NodeType == IMPORT_TYPE {
						var import_string string
						for i := 0; i < len(lineTokens); i++ {
							if i == len(lineTokens) - 1 {
								// because this is the path like for example "./some/path/App.js"
								// we need to check the valid path
								import_string += "\"" + lineTokens[i] + "\""
							}else {
								import_string += lineTokens[i] + " ";
							}							
						}
						topNode.Imports = append(topNode.Imports, import_string)
					} else if topNode.NodeType == ROUTER_TYPE {
						var route_string string
						for i := 0; i < len(lineTokens); i++ {
							if i == len(lineTokens) - 1 {
								route_string += lineTokens[i]
							} else {
								route_string += lineTokens[i] + " "
							}
						}
						topNode.Router = append(topNode.Router, route_string)
					} else {
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
				// fmt.Println(topNode.DirectInfo, topNode.NodeType)
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
		return stack.Peek(), imports, routes
	}

	return stack.Peek(), imports, routes
}

func checkToken(token string) int {
	if Keywords[token] {
		return KEYWORD_TYPE
	}else if config[token] {
		return CONFIG_TYPE
	}else if import_s[token] {
		return IMPORT_TYPE
	} else if router_s[token] {
		return ROUTER_TYPE
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
