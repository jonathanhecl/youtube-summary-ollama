package main

import (
	"fmt"
)

var (
	version = "0.0.1"
)

func main() {
	fmt.Println("YouTube Summary Ollama v", version)

	videoID := "TNrylG4hUIQ"

	content, err := getTranscription(videoID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content)

}
