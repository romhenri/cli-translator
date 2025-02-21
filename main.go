// CLI Translator by @romhenri
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"cli-translator/services"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("  Single: cli-translator <text> [-lang] [-from:source] [-d]")
		fmt.Println("  Continuous: cli-translator -lang [-from:source]")
		return
	}

	// Default values
	var text string
	fromLang := "auto"
	targetLang := "en"
	includeDetails := false

	if len(os.Args) > 1 {
		if strings.HasPrefix(os.Args[1], "-") && !strings.HasPrefix(os.Args[1], "-from:") {
			targetLang = strings.TrimPrefix(os.Args[1], "-")
			
			// "-from:" in Continuous Mode
			for _, arg := range os.Args[2:] {
				if strings.HasPrefix(arg, "-from:") {
					fromLang = strings.TrimPrefix(arg, "-from:")
					if fromLang == "" {
						fromLang = "auto"
					}
				}
			}

			fmt.Printf("CLI-Translator [to %s] [from %s]\n", targetLang, fromLang)
			continuousMode(targetLang, fromLang, includeDetails)
			return
		}
		text = os.Args[1]
	}

	if len(os.Args) > 2 {
		for _, arg := range os.Args[2:] {
			if strings.HasPrefix(arg, "-from:") {
				fromLang = strings.TrimPrefix(arg, "-from:")
				if fromLang == "" {
					fromLang = "auto"
				}
			} else if arg == "-d" {
				includeDetails = true
			} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
				targetLang = strings.TrimPrefix(arg, "-")
			} else {
				if text != "" {
					text += " " + arg
				} else {
					text = arg
				}
			}
		}
	}

	// Single Mode
	translation, err := services.Translate(text, targetLang, fromLang, includeDetails)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(">", translation)
}

// Continuous Mode
func continuousMode(targetLang, fromLang string, includeDetails bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if text == "0" {
			fmt.Println("Exiting...")
			break
		}

		translation, err := services.Translate(text, targetLang, fromLang, includeDetails)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println(translation, "\n")
	}
}
