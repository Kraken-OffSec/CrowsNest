package cmd

import (
	"crowsnest/internal/easyTime"
	"crowsnest/internal/pretty"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	// Add Subcommand to db command
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().IntVarP(&logLastXLogs, "last", "l", 0, "Number of logs to show")
	logsCmd.Flags().StringVarP(&logStartDate, "start", "s", "", "Start date for logs")
	logsCmd.Flags().StringVarP(&logEndDate, "end", "e", "", "End date for logs")
	logsCmd.Flags().StringVarP(&logSeverity, "severity", "v", "", "Comma delimited Log severity to show (info, error, fatal)")
}

// LogEntry represents a parsed log entry
type LogEntry struct {
	Level      string                 `json:"level"`
	Timestamp  string                 `json:"ts"` // Store as string initially
	Message    string                 `json:"msg"`
	Details    map[string]interface{} `json:"-"`
	ParsedTime time.Time              // Parsed time for sorting and filtering
}

var (
	logLastXLogs int
	logStartDate string
	logEndDate   string
	logSeverity  string

	logsCmd = &cobra.Command{
		Use:   "logs",
		Short: "View logs",
		Long:  `View logs for the Dehasher CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			logsPath := filepath.Join(os.Getenv("HOME"), ".local", "share", "Dehasher", "logs")

			if logSeverity == "" {
				logSeverity = "info,error,fatal"
			} else {
				logSeverity = strings.ToLower(logSeverity)
			}

			var (
				logInfo  bool
				logError bool
				logFatal bool
			)
			if strings.Contains(logSeverity, "info") {
				logInfo = true
			}
			if strings.Contains(logSeverity, "error") {
				logError = true
			}
			if strings.Contains(logSeverity, "fatal") {
				logFatal = true
			}

			var allLogs []string
			if logSeverity == "info" {
				allLogs = append(allLogs, filepath.Join(logsPath, "info.log"))
			} else if logSeverity == "error" || logSeverity == "fatal" {
				allLogs = append(allLogs, filepath.Join(logsPath, "error.log"))
			} else {
				allLogs = append(allLogs, filepath.Join(logsPath, "info.log"), filepath.Join(logsPath, "error.log"))
			}

			var timeChunk easyTime.TimeChunk
			if logStartDate != "" {
				timeChunk = easyTime.NewTimeChunk(logStartDate, logEndDate, debugGlobal)
			}

			var parsedLogs []LogEntry
			for _, logFile := range allLogs {
				// Read the log file
				logData, err := os.ReadFile(logFile)
				if err != nil {
					fmt.Printf("Error reading log file %s: %v\n", logFile, err)
					continue
				}

				// Split the file into lines
				logLines := strings.Split(strings.TrimSpace(string(logData)), "\n")
				for _, line := range logLines {
					if line == "" {
						continue
					}

					// Parse the JSON log entry
					var entry LogEntry
					var rawEntry map[string]interface{}
					if err := json.Unmarshal([]byte(line), &entry); err != nil {
						fmt.Printf("Error parsing log entry: %v\n", err)
						continue
					}

					// Unmarshal to get additional fields
					if err := json.Unmarshal([]byte(line), &rawEntry); err != nil {
						fmt.Printf("Error parsing raw log entry: %v\n", err)
						continue
					}

					// Parse the timestamp
					parsedTime, err := time.Parse("2006-01-02T15:04:05.999-0700", entry.Timestamp)
					if err != nil {
						// Try RFC3339
						parsedTime, err = time.Parse(time.RFC3339, entry.Timestamp)
						if err != nil {
							// Try RFC3339Nano
							parsedTime, err = time.Parse(time.RFC3339Nano, entry.Timestamp)
							if err != nil {
								fmt.Printf("Error parsing timestamp '%s': %v\n", entry.Timestamp, err)
								continue
							}
						}
					}
					entry.ParsedTime = parsedTime

					// Store additional fields in Details
					entry.Details = make(map[string]interface{})
					for k, v := range rawEntry {
						if k != "level" && k != "ts" && k != "msg" {
							entry.Details[k] = v
						}
					}

					// Filter by severity
					if (logInfo && strings.EqualFold(entry.Level, "INFO")) ||
						(logError && strings.EqualFold(entry.Level, "ERROR")) ||
						(logFatal && strings.EqualFold(entry.Level, "FATAL")) {

						// Filter by date range if specified
						if timeChunk.IsSet() {
							if entry.ParsedTime.Before(timeChunk.StartTime) || entry.ParsedTime.After(timeChunk.EndTime) {
								continue
							}
						}

						parsedLogs = append(parsedLogs, entry)
					}
				}
			}

			// Limit the number of logs if specified
			if logLastXLogs > 0 && len(parsedLogs) > logLastXLogs {
				// Sort logs by timestamp (newest first)
				// This is a simple bubble sort for demonstration
				for i := 0; i < len(parsedLogs)-1; i++ {
					for j := 0; j < len(parsedLogs)-i-1; j++ {
						if parsedLogs[j].ParsedTime.Before(parsedLogs[j+1].ParsedTime) {
							parsedLogs[j], parsedLogs[j+1] = parsedLogs[j+1], parsedLogs[j]
						}
					}
				}
				parsedLogs = parsedLogs[:logLastXLogs]
			}

			// Display logs in a table
			if len(parsedLogs) > 0 {
				// Prepare table headers and rows
				headers := []string{"Date", "Severity", "Message", "Details"}
				rows := make([][]string, len(parsedLogs))

				for i, log := range parsedLogs {
					// Format timestamp
					timeStr := log.ParsedTime.Format("2006-01-02 15:04:05")

					// Format details
					detailsStr := ""
					for k, v := range log.Details {
						detailsStr += fmt.Sprintf("%s: %v, ", k, v)
					}
					if len(detailsStr) > 2 {
						detailsStr = detailsStr[:len(detailsStr)-2] // Remove trailing comma and space
					}

					rows[i] = []string{timeStr, log.Level, log.Message, detailsStr}
				}

				// Display the table
				pretty.Table(headers, rows)
			} else {
				fmt.Println("No logs found matching the specified criteria.")
			}
		},
	}
)

type Severity int

const (
	INFO Severity = iota
	ERROR
	FATAL
	UNKNOWN Severity = -1
)
