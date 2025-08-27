package cmdline

import (
	"fmt"

	"github.com/nikhilbhatia08/eml/cmd/executor"
	"github.com/nikhilbhatia08/eml/parser"
	"github.com/nikhilbhatia08/eml/codegen"
)

func ParseCommandLineArgs(args []string) {
	command := args[0]
	if command == "eml" {
		arg1 := args[1]
		if arg1 == "make" {
			arg2 := args[2]
			op := executor.ExecuteMakeCommand(arg2)
			if op == "Success" {
				// Handle success
				fmt.Println("Success")
			}else {
				fmt.Println("Error")
			}
		}else if arg1 == "run" {
			executor.RunProject()
		}else if arg1 == "compile" {
			// parsing
			root := parser.Parser()
			// code generation
			lines := codegen.GenerateHtmlCodeFromAST(root)
			for _, line := range lines {
				fmt.Println(line)
			}
		}
	}
}