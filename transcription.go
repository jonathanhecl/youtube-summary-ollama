package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Caption struct {
	Start    float64
	Duration float64
	Text     string
}

type Transcription struct {
	Caption []Caption
}

type transcript struct {
	XMLName xml.Name `xml:"transcript"`
	Text    []struct {
		XMLName xml.Name `xml:"text"`
		Start   string   `xml:"start,attr"`
		Dur     string   `xml:"dur,attr"`
		Text    string   `xml:",chardata"`
	} `xml:"text"`
}

func getTranscription(videoID string) (*Transcription, string, error) {
	var description string
	content, err := _getHTML("https://www.youtube.com/watch?v=" + videoID)
	if err != nil {
		return nil, description, err
	}

	var re1 = regexp.MustCompile(`(?m)"attributedDescription":{"content":"([^"]*)","commandRuns":\[{"startIndex"`)
	matches := re1.FindStringSubmatch(content)

	if len(matches) == 2 {
		description = matches[1]
	}

	var re2 = regexp.MustCompile(`(?m){"captionTracks":\[{"baseUrl":"([a-zA-Z0-9\\:.,_\-/?=]*)","name"`)
	matches = re2.FindStringSubmatch(content)

	if len(matches) < 2 {
		fmt.Println(matches)
		return nil, description, errors.New("No matches found")
	}

	urlCaptions := matches[1]
	urlCaptions = strings.Replace(urlCaptions, "\\u0026", "&", -1)

	content, err = _getHTML(urlCaptions)
	if err != nil {
		return nil, description, err
	}

	var t transcript
	err = xml.Unmarshal([]byte(content), &t)
	if err != nil {
		return nil, description, err
	}

	captions := make([]Caption, len(t.Text))
	for i, text := range t.Text {
		start, _ := strconv.ParseFloat(text.Start, 64)
		dur, _ := strconv.ParseFloat(text.Dur, 64)

		captions[i] = Caption{
			Start:    start,
			Duration: dur,
			Text:     text.Text,
		}
	}

	return &Transcription{Caption: captions}, description, nil
}

func (t *Transcription) GetText() string {
	text := ""

	for _, caption := range t.Caption {
		text += caption.Text + " "
	}

	return text
}
