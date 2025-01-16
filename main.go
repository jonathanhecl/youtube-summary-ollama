package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jonathanhecl/chunker"
	"github.com/jonathanhecl/gollama"
)

var (
	version   string = "0.0.1"
	maxLength int64  = 1024 * 128
)

func main() {
	fmt.Println("YouTube Summary Ollama v", version)

	if len(os.Args) < 2 {
		fmt.Println("Usage: youtube-summary-ollama <video_id>")
		return
	}

	m := gollama.New("command-r7b")
	m.SetTemperature(1)
	m.SetContextLength(maxLength)

	videoID := os.Args[1]

	content, desc, err := getTranscription(videoID)
	if err != nil {
		log.Fatal(err)
	}

	data := content.GetText()

	fmt.Println("Extracting Description... (length:", len(desc), ")")

	res, err := m.Chat(context.Background(), `
	Extract the main topic (with details) from the following description:
	"`+desc+`"

	Avoid information about social media and disclaimers.
	Avoid interpretations.
	If the description is empty, respond with "No description".
	`)

	if err != nil {
		log.Fatal(err)
	}

	desc = res.Content

	fmt.Println("Main Topics:", res.Content)

	fmt.Println("Sending to Ollama... (length:", len(data), ")")

	resume := ""
	if len(data) < int(maxLength) {

		res1, err := m.Chat(context.Background(), `
		Create a summary of the topics discussed in this video transcript:
		"`+data+`"

		This is the short description of the video:
		"`+desc+`"

		SUMMARY:`)
		if err != nil {
			log.Fatal(err)
		}

		resume = res1.Content
	} else {
		chunk := chunker.NewChunker(int(maxLength), 128, chunker.DefaultSeparators, true, false)

		for _, chunk := range chunk.Chunk(data) {

			res1, err := m.Chat(context.Background(), `
		Create a summary of the topics discussed in this video transcript:
		"`+chunk+`"

		This is the short description of the video:
		"`+desc+`"

		SUMMARY:`)
			if err != nil {
				log.Fatal(err)
			}

			resume += res1.Content
		}
	}

	fmt.Println("Resume:", resume)
	fmt.Println("Reinterpretando... (length:", len(resume), ")")

	type ResTopics struct {
		Topics []string `json:"topics"`
	}

	prompt := `Make a synthesis in clear points of the previous summary in Spanish:
"` + resume + `"

- Don't do a literal translation.
- Only the video topics (up to 3 words).
- No conclusions or analysis.
- Respond in JSON`

	r, _ := gollama.StructToStructuredFormat(ResTopics{})

	res2, err := m.Chat(context.Background(), prompt, r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res2)
}
