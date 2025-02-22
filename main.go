// CLI Translator by @romhenri
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"cli-translator/services"
	"cli-translator/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("  Single: cli-translator <text> [-lang] [-from:source] [-d] [-debug]")
		fmt.Println("  Continuous: cli-translator -lang [-from:source] [-debug]")
		return
	}

	// Version flag
	if os.Args[1] == "-version" || os.Args[1] == "-v" {
		fmt.Println("CLI Translator", config.Version)
		return
	}

	// Default values
	var text string
	fromLang := "auto"
	targetLang := "en"
	includeDetails := false
	debugMode := false

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
				} else if arg == "-debug" {
					debugMode = true
				}
			}

			fmt.Printf("CLI-Translator [to %s] [from %s]\n", targetLang, fromLang)
			continuousMode(targetLang, fromLang, includeDetails, debugMode)
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
			} else if arg == "-debug" {
				debugMode = true
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

	if debugMode {
		fmt.Println("//")
		fmt.Println("Debug: Text to translate:", text)
		fmt.Println("Debug: Target Language:", targetLang)
		fmt.Println("Debug: Source Language:", fromLang)
		fmt.Println("Debug: Include Details:", includeDetails)
		fmt.Println("//\n")
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
func continuousMode(targetLang, fromLang string, includeDetails, debugMode bool) {
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

		if debugMode {
			fmt.Println("//")
			fmt.Println("Debug: Text to translate:", text)
			fmt.Println("Debug: Target Language:", targetLang)
			fmt.Println("Debug: Source Language:", fromLang)
			fmt.Println("Debug: Include Details:", includeDetails)
			fmt.Println("//\n")
		}

		translation, err := services.Translate(text, targetLang, fromLang, includeDetails)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println(translation, "\n")
	}
}
