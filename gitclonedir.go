package main

import (
	"os"
)

func main() {
	githubUrl := os.Args[1]
	outputPath := "."

	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	target, _ := InferTargetFromUrl(githubUrl)

	downloadCommand := NewDownloadCommand(target, outputPath)
	err := downloadCommand.Execute()

	if err != nil {
		println("Error while downloading repository content:", err.Error())
		return
	}
}
