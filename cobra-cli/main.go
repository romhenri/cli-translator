package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"cli-translator/services"
	"cli-translator/handlers"
	"cli-translator/config"

	"github.com/spf13/cobra"
)

var (
	fromLang      string
	targetLang    string
	includeDetails bool
	debugMode     bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli-translator",
		Short: "A simple CLI translator",
	}

	translateCmd := &cobra.Command{
		Use:   "tl [text]",
		Short: "Translate a given text",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			text := strings.Join(args, " ")
			translation, err := services.Translate(text, targetLang, fromLang, includeDetails)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			history.SaveToHistory(text, fromLang, targetLang, translation)
			fmt.Println(">", translation)
		},
	}

	continuousCmd := &cobra.Command{
		Use:   "cont",
		Short: "Enter continuous translation mode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("CLI-Translator [to %s] [from %s]\n", targetLang, fromLang)
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
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the CLI Translator version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.CobraAppName, config.CobraVersion)
		},
	}

	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show translation history",
		Run: func(cmd *cobra.Command, args []string) {
			history.ShowHistory()
		},
	}

	clearHistoryCmd := &cobra.Command{
		Use:   "clear-history",
		Short: "Clear translation history",
		Run: func(cmd *cobra.Command, args []string) {
			history.ClearHistory()
		},
	}

	// Flags
	translateCmd.Flags().StringVarP(&fromLang, "from", "f", "auto", "Source language")
	translateCmd.Flags().StringVarP(&targetLang, "to", "t", "en", "Target language")
	translateCmd.Flags().BoolVarP(&includeDetails, "details", "d", false, "Include translation details")
	translateCmd.Flags().BoolVarP(&debugMode, "debug", "D", false, "Enable debug mode")

	continuousCmd.Flags().StringVarP(&fromLang, "from", "f", "auto", "Source language")
	continuousCmd.Flags().StringVarP(&targetLang, "to", "t", "en", "Target language")
	continuousCmd.Flags().BoolVarP(&includeDetails, "details", "d", false, "Include translation details")
	continuousCmd.Flags().BoolVarP(&debugMode, "debug", "D", false, "Enable debug mode")

	rootCmd.AddCommand(translateCmd, continuousCmd, versionCmd, historyCmd, clearHistoryCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}