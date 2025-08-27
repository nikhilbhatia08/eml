package main

import (
	"fmt"

	"github.com/nikhilbhatia08/eml/cmd/cmdline"
	"github.com/nikhilbhatia08/eml/cmd/utils"
)

func main() {
	arguementArray := utils.GetCommandLineArgs()
	fmt.Println("Arg Array: ", arguementArray)
	cmdline.ParseCommandLineArgs(arguementArray)
}