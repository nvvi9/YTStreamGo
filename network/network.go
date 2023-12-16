package network

import (
	"fmt"
	"io"
	"net/http"
)

var client = &http.Client{}

func GetVideoPage(id string) (string, error) {
	return getRaw(fmt.Sprintf("https://www.youtube.com/watch?v=%s", id))
}

func GetJsFile(jsPath string) (string, error) {
	return getRaw(fmt.Sprintf("https://www.youtube.com%s", jsPath))
}

func getRaw(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.98 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
