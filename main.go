package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var DEBUG_MODE bool = false

func translate(text string, target string, from string) (string, error) {
	apiURL := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s", from, target, url.QueryEscape(text))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("error on request make: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on request exec: %v", err)
	}
	defer resp.Body.Close()

	// Response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error: status %d, response: %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error on response read: %v", err)
	}

	if DEBUG_MODE {
		fmt.Println("URL:", apiURL)
		fmt.Println("JSON recebido:", string(bodyBytes))
	}

	var result []interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("error on JSON read: %v", err)
	}

	if len(result) > 0 {
		if firstArray, ok := result[0].([]interface{}); ok && len(firstArray) > 0 {
			if translationData, ok := firstArray[0].([]interface{}); ok && len(translationData) > 0 {
				if translatedText, ok := translationData[0].(string); ok {
					return translatedText, nil
				}
			}
		}
	}

	return "", fmt.Errorf("not found")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: cli-translater <texto> [-idioma]")
		return
	}

	text := os.Args[1]
	fromLang := "auto"
	targetLang := "en"

	if len(os.Args) > 2 {
		if strings.HasPrefix(os.Args[2], "-") {
			targetLang = strings.TrimPrefix(os.Args[2], "-")
		}

		if len(os.Args) > 3 {
			if strings.HasPrefix(os.Args[3], "-from:") {
				fromLang = strings.TrimPrefix(os.Args[3], "-from:")
				if fromLang == "" {
					fromLang = "auto"
				}
			}
		}
	}

	translation, err := translate(text, targetLang, fromLang)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(">", translation)
}
