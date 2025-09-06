package executor

import (
	// "fmt"
	"fmt"
	"os"
	"os/exec"
)

func ExecuteMakeCommand(nameOfProject string) string {
	// We need to get rid of this SomeHow
	dir := nameOfProject
	err := os.Mkdir(dir, 0755)
	if err != nil {
		fmt.Println("Error creating directory : ", err)
		return "Error"
	}
	err = os.Chdir(dir)
	if err != nil {
		fmt.Println("Error changing Directory : ", err)
		return "Error"
	}

	WriteGitIgnore()

	WriteMain()

	_, err = exec.Command("npx", "create-react-app", "build").Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "Error"
	}
	os.Chdir("build")
	// Generating Tailwindcss
	_, err = exec.Command("npm", "install", "-D", "tailwindcss@3").Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "Error"
	}

	_, err = exec.Command("npm", "install", "react-router").Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "Error"
	}

	os.Remove("tailwind.config.js")
	// Creating Tailwind config file
	WriteToTailwindConfig()

	os.Chdir("src")
	// Removing index.css to be rewritten
	os.Remove("index.css")
	WriteToIndexCss()

	os.Remove("index.js")
	WriteToIndexjs()
	// Removing App.js To be Rewritten
	os.Remove("App.js")
	WriteToFile()
	return "Success"
}

func WriteToFile() {
	file, err := os.Create("App.js") // truncates if file exists
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(WriteAppjs())
	if err != nil {
		panic(err)
	}

	fmt.Println("File written successfully")
}

func WriteAppjs() string {
	return "import React from 'react';\n\nfunction App() {\n\treturn (\n\t\t<div>\n\t\t\t<h1 className=\"text-3xl font-bold underline\">Hello, World!</h1>\n\t\t</div>\n\t);\n}\n\nexport default App;\n"
}

func WriteMain() {
	file, err := os.Create("main.ehtml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("div:\n\th1: Hello World\n")
	if err != nil {
		panic(err)
	}

	fmt.Println("File written successfully")
}

func WriteToTailwindConfig() {
	file, err := os.Create("tailwind.config.js")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("module.exports = {\n\tcontent: [\n\t\t'./src/**/*.{js,jsx,ts,tsx}',\n\t],\n\ttheme: {\n\t\textend: {},\n\t},\n\tplugins: [],\n};\n")
	if err != nil {
		panic(err)
	}

	fmt.Println("Tailwind config file created successfully")
}

func WriteToIndexjs(){
	file, err := os.Create("index.js")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("import React from 'react';\nimport './index.css'\nimport ReactDOM from 'react-dom/client';\nimport { BrowserRouter } from 'react-router';\nimport App from './App';\n\nconst root = document.getElementById('root');\n\nReactDOM.createRoot(root).render(\n\t<BrowserRouter>\n\t\t<App />\n\t</BrowserRouter>\n);")
	if err != nil {
		panic(err)
	}

	fmt.Println("Index js file created successfully")
}

func WriteToIndexCss() {
	file, err := os.Create("index.css")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n")
	if err != nil {
		panic(err)
	}

	fmt.Println("Index css file created successfully")
}

func WriteGitIgnore() {
	file, err := os.Create(".gitignore")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("build/\n")
	if err != nil {
		panic(err)
	}

	fmt.Println("Git ignore file created successfully")
}