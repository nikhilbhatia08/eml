package utils

import "os"

func GetCommandLineArgs() []string {
	commandLineArgs := os.Args[1:]
	return commandLineArgs
}