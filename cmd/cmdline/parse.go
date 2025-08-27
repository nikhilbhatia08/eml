package cmdline

import (
	"fmt"
	"os"

	"github.com/nikhilbhatia08/eml/cmd/executor"
	"github.com/nikhilbhatia08/eml/codegen"
	"github.com/nikhilbhatia08/eml/parser"
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
			targets := parser.Parser()
			// code generation
			for _, target := range targets {
				fmt.Println(target.Path)
				fmt.Println("Imports: ")
				imports := codegen.GenerateImports(target.Imports)
				for _, import_s := range imports {
					fmt.Println(import_s)
				}
				lines := codegen.GenerateHtmlCodeFromAST(target.Root)
				for _, line := range lines {
					fmt.Println(line)
				}
				os.Chdir("build")
				os.Chdir("src")
				err := codegen.WriteToFile(target.Path, lines, imports)
				if err != nil {
					fmt.Println("Error writing")
				}
			}

		}
	}
}