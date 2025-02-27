package logic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"time"
)

type Logs struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
}

func GetLogs(filename string, logs chan string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logs <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(logs)
}

func SortByLevel(logSlice []string, level string, ch chan string) {
	pattern := fmt.Sprintf(`.*%s.*`, level)
	re := regexp.MustCompile(pattern)
	for _, lg := range logSlice {
		if re.FindString(lg) != "" {
			ch <- lg
		}
	}
	close(ch)
}

func SortByTime(logSlice []string, since string, until string, ch chan string) {
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})`)
	var logEntries []Logs

	for _, log := range logSlice {
		match := re.FindString(log)
		if match != "" {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", match)
			if err == nil {
				logEntries = append(logEntries, Logs{
					Time:    parsedTime,
					Message: log,
				})
			}
		}
	}

	var filteredLogs []Logs
	for _, entry := range logEntries {
		if since != "" {
			sinceTime, err := time.Parse("2006-01-02 15:04:05", since)
			if err == nil && entry.Time.Before(sinceTime) {
				continue
			}
		}
		if until != "" {
			untilTime, err := time.Parse("2006-01-02 15:04:05", until)
			if err == nil && entry.Time.After(untilTime) {
				continue
			}
		}

		filteredLogs = append(filteredLogs, entry)
	}

	sort.Slice(filteredLogs, func(i, j int) bool {
		return filteredLogs[i].Time.Before(filteredLogs[j].Time)
	})

	for _, entry := range filteredLogs {
		ch <- entry.Message
	}
	close(ch)
}

func WriteLogs(logs []string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, lg := range logs {
		_, err := w.WriteString(lg + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func SaveJson(logs []string, fileName string) error {
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+(INFO|ERROR|WARN|DEBUG)\s+(.*)$`)

	var logEntries []Logs

	for _, logStr := range logs {
		matches := re.FindStringSubmatch(logStr)
		if len(matches) == 4 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", matches[1])
			if err != nil {
				return err
			}
			logEntry := Logs{
				Time:    parsedTime,
				Level:   matches[2],
				Message: matches[3],
			}
			logEntries = append(logEntries, logEntry)
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(logEntries)
}
