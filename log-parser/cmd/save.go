/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log-parser/logic/log"
	"strings"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save the parsed logs to a file",
	Long: `The 'save' command allows you to save the filtered or parsed logs 
to a specified file. You can choose between different formats such as plain text or JSON.

Example usage:
logparser save --file /path/to/output.txt --format text
logparser save --file /path/to/output.json --format json`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		format, _ := cmd.Flags().GetString("format")
		file = strings.TrimSpace(file)
		ch := make(chan string)
		logs := make([]string, 0)
		go logic.GetLogs("temp.txt", ch)
		for lg := range ch {
			logs = append(logs, lg)
		}

		if format == "text" {
			logic.WriteLogs(logs, file)
			fmt.Println("Logs saved at %s", file)
		} else {
			err := logic.SaveJson(logs, file)
			if err != nil {
				fmt.Printf("Error saving logs in JSON format: %v\n", err)
			} else {
				fmt.Printf("Logs saved in JSON format at %s\n", file)
			}
		}

	},
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "p", "text", "saving format")
	rootCmd.AddCommand(saveCmd)
}
