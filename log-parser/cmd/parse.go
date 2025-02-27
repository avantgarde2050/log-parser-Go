package cmd

import (
	"errors"
	"fmt"
	"log"
	"log-parser/logic/log"
	"strings"

	"github.com/spf13/cobra"
)

var parsingCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse log files and display extracted information",
	Long: `The 'parsing' command allows you to parse log files, extracting useful data such as 
the event timestamp, event level (INFO, ERROR, etc.), and the event description.
You can filter the logs by level or date range for more specific results.

Example usage:

logparser parse --file /path/to/logfile.log --level ERROR
logparser parse --file /path/to/logfile.log --since 2025-02-14 00:00:00 --until 2025-02-15 23:59:59
logparser parse --file /path/to/logfile.log --level INFO --since 2025-02-14 00:00:00 --until 2025-02-15 23:59:59`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		level, _ := cmd.Flags().GetString("level")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")
		file = strings.TrimSpace(file)
		level = strings.TrimSpace(level)
		since = strings.TrimSpace(since)
		until = strings.TrimSpace(until)
		if file == "" {

			log.Fatalf("%v", errors.New("You must have to write filename"))
		}

		fmt.Printf("Parsing log file: %s\n", file)
		if level != "" {
			fmt.Printf("Level: %s\n", level)
		}
		if since != "" && until != "" {
			fmt.Printf("Filtering by time: from %s to %s\n", since, until)
		}
		var logSlice = make([]string, 0)
		var logs chan string = make(chan string)
		go logic.GetLogs(file, logs)
		for lg := range logs {
			logSlice = append(logSlice, lg)
		}
		var logsLevelSorted = make(chan string)
		go logic.SortByLevel(logSlice, level, logsLevelSorted)
		var levelSortedLogs = make([]string, 0)
		for lg := range logsLevelSorted {
			levelSortedLogs = append(levelSortedLogs, lg)
		}
		var finalLogs = make(chan string)
		var tempLogs = make([]string, 0)
		go logic.SortByTime(levelSortedLogs, since, until, finalLogs)
		for lg := range finalLogs {
			fmt.Println(lg)
			tempLogs = append(tempLogs, lg)
		}

		logic.WriteLogs(logSlice, "logs.txt")
		logic.WriteLogs(tempLogs, "temp.txt")
	},
}

func init() {

	rootCmd.PersistentFlags().StringP("file", "f", "logs.txt", "Path to the log file to be parsed")
	rootCmd.PersistentFlags().StringP("level", "l", "", "Filter logs by severity level (e.g., ERROR, INFO)")
	rootCmd.PersistentFlags().StringP("since", "s", "", "Filter logs from a specific time (format: YYYY-MM-DD HH:MM:SS)")
	rootCmd.PersistentFlags().StringP("until", "u", "", "Filter logs until a specific time (format: YYYY-MM-DD HH:MM:SS)")
	rootCmd.AddCommand(parsingCmd)
	rootCmd.AddCommand(filterCmd)
	rootCmd.AddCommand(saveCmd)
}
