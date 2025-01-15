package main

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func getTranscription(videoID string) (string, error) {
	content, err := getHTML("https://www.youtube.com/watch?v=" + videoID)
	if err != nil {
		return "", err
	}

	var re = regexp.MustCompile(`(?m){"captionTracks":\[{"baseUrl":"([a-zA-Z0-9\\:.,\-/?=]*)","name"`)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return "", errors.New("No matches found")
	}

	urlCaptions := matches[1]
	urlCaptions = strings.Replace(urlCaptions, "\\u0026", "&", -1)

	content, err = getHTML(urlCaptions)
	if err != nil {
		return "", err
	}

	return content, nil
}

func getHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
