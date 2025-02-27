package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log-parser/logic/log"
	"strings"
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter logs by severity level and/or time range",
	Long: `The 'filter' command allows you to filter log entries based on their severity level (e.g., INFO, ERROR) 
and within a specified time range. This command provides an efficient way to narrow down the logs you're interested in 
and display them in a readable format.

Examples of usage:

1. Filter logs by severity level:
   logparser filter --level ERROR

2. Filter logs by time range:
   logparser filter --since "2025-02-14 00:00:00" --until "2025-02-15 23:59:59"

3. Filter logs by both severity level and time range:
   logparser filter --level INFO --since "2025-02-14 00:00:00" --until "2025-02-15 23:59:59"`,
	Run: func(cmd *cobra.Command, args []string) {
		level, _ := cmd.Flags().GetString("level")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")
		level = strings.TrimSpace(level)
		since = strings.TrimSpace(since)
		until = strings.TrimSpace(until)
		if level != "" {
			fmt.Printf("Level: %s\n", level)
		}
		if since != "" && until != "" {
			fmt.Printf("Filtering by time: from %s to %s\n", since, until)
		}
		var logSlice = make([]string, 0)
		var logs chan string = make(chan string)
		go logic.GetLogs("logs.txt", logs)
		for lg := range logs {
			logSlice = append(logSlice, lg)
		}
		var logsLevelSorted = make(chan string)
		go logic.SortByLevel(logSlice, level, logsLevelSorted)
		var levelSortedLogs = make([]string, 0)
		for lg := range logsLevelSorted {
			levelSortedLogs = append(levelSortedLogs, lg)
		}
		var endLogs = make([]string, 0)
		var finalLogs = make(chan string)
		go logic.SortByTime(levelSortedLogs, since, until, finalLogs)
		for lg := range finalLogs {
			fmt.Println(lg)
			endLogs = append(endLogs, lg)
		}
		logic.WriteLogs(endLogs, "temp.txt")
	},
}

func init() {

}
