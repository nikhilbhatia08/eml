package executor

import (
	"fmt"
	"os"
	"os/exec"
)

func RunProject() {
	os.Chdir("build")
	op, _ := exec.Command("pwd").Output()
	fmt.Println(string(op))
	_, err := exec.Command("npm", "start").Output()
	if err != nil {
		// Handle error
		fmt.Println("Error executing command:", err)
		return
	}
}