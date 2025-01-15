package main

import (
	"encoding/xml"
	"errors"
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

func getTranscription(videoID string) (*Transcription, error) {
	content, err := _getHTML("https://www.youtube.com/watch?v=" + videoID)
	if err != nil {
		return nil, err
	}

	var re = regexp.MustCompile(`(?m){"captionTracks":\[{"baseUrl":"([a-zA-Z0-9\\:.,\-/?=]*)","name"`)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return nil, errors.New("No matches found")
	}

	urlCaptions := matches[1]
	urlCaptions = strings.Replace(urlCaptions, "\\u0026", "&", -1)

	content, err = _getHTML(urlCaptions)
	if err != nil {
		return nil, err
	}

	var t transcript
	err = xml.Unmarshal([]byte(content), &t)
	if err != nil {
		return nil, err
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

	return &Transcription{Caption: captions}, nil
}
