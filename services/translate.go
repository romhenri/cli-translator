package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var DEBUG_MODE bool = false

func Translate(text string, target string, from string, includeDetails bool) (string, error) {
	detailParams := "t"
	if includeDetails {
		detailParams = "t&dt=bd"
	}

	apiURL := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=%s&q=%s",
		from, target, detailParams, url.QueryEscape(text))

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

	var translatedText string
	var synonyms []string

	if len(result) > 0 {
		// Translation
		if firstArray, ok := result[0].([]interface{}); ok && len(firstArray) > 0 {
			if translationData, ok := firstArray[0].([]interface{}); ok && len(translationData) > 0 {
				if text, ok := translationData[0].(string); ok {
					translatedText = text
				}
			}
		}

		// Synonyms
		if includeDetails && len(result) > 1 {
			if definitions, ok := result[1].([]interface{}); ok {
				for _, def := range definitions {
					if defArray, ok := def.([]interface{}); ok && len(defArray) > 1 {
						grammaticalType := defArray[0].(string)
						if synonymsList, ok := defArray[1].([]interface{}); ok {
							for _, synonym := range synonymsList {
								if synonymStr, ok := synonym.(string); ok {
									capitalizedSynonym := strings.ToUpper(synonymStr[:1]) + synonymStr[1:]
									synonyms = append(synonyms, fmt.Sprintf("%s (%s)", capitalizedSynonym, grammaticalType))
								}
							}
						}
					}
				}
			}
		}
	}

	if translatedText == "" {
		return "", fmt.Errorf("translation not found")
	}

	if includeDetails && len(synonyms) > 0 {
		return fmt.Sprintf("%s\n\n- %s", translatedText, strings.Join(synonyms, "\n- ")), nil
	}

	return translatedText, nil
}
