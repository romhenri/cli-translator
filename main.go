// CLI Translator by @romhenri
package main

import (
	"fmt"
	"os"
	"strings"
	"cli-translator/services"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: cli-translater <texto> [-idioma] [-d]")
		return
	}

	// Default values
	text := os.Args[1]
	fromLang := "auto"
	targetLang := "en"
	includeDetails := false

	if len(os.Args) > 2 {
		if strings.HasPrefix(os.Args[2], "-") {
			targetLang = strings.TrimPrefix(os.Args[2], "-")
		}

		if len(os.Args) > 3 {
			for _, arg := range os.Args[3:] {
				if strings.HasPrefix(arg, "-from:") {
					fromLang = strings.TrimPrefix(arg, "-from:")
					if fromLang == "" {
						fromLang = "auto"
					}
				} else if arg == "-d" {
					includeDetails = true
				}
			}
		}
	}
	translation, err := services.Translate(text, targetLang, fromLang, includeDetails)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(">", translation)
}
