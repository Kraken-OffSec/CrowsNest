package easyTime

import (
	"crowsnest/internal/debug"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

type TimeChunk struct {
	StartTime time.Time
	EndTime   time.Time
	set       bool
}

func (tc *TimeChunk) isValid() bool {
	if !tc.StartTime.IsZero() && !tc.EndTime.IsZero() && tc.StartTime.Before(tc.EndTime) {
		tc.set = true
		return true
	}
	tc.set = false
	return false
}

func (tc *TimeChunk) IsSet() bool {
	return tc.set
}

func NewTimeChunk(start, end string, debugOn bool) TimeChunk {
	if debugOn {
		debug.PrintInfo("parsing time chunk")
		debug.PrintInfo(fmt.Sprintf("Start: %s, End: %s", start, end))
		zap.L().Info("parsing time chunk",
			zap.String("start", start),
			zap.String("end", end),
		)
	}

	if end == "" {
		if debugOn {
			debug.PrintInfo("no end time provided, using now")
		}
		end = "now"
	}

	tc := TimeChunk{
		StartTime: parseUserTime(start),
		EndTime:   parseUserTime(end),
	}

	if debugOn {
		debug.PrintInfo("checking if time chunk is valid")
		debug.PrintInfo(fmt.Sprintf("Start: %s, End: %s", tc.StartTime, tc.EndTime))
	}
	if !tc.isValid() {
		fmt.Println("[!] Invalid time chunk")
		zap.L().Fatal("invalid_time_chunk",
			zap.String("message", "invalid time chunk"),
		)
		os.Exit(1)
	}

	return tc
}

func parseUserTime(args string) time.Time {
	args = strings.TrimSpace(args)

	if strings.EqualFold(args, "now") {
		return time.Now()
	}

	// Check if time contains a space, if so, it's in 'last 24 hours' format
	if strings.Contains(args, " ") && !containsMonth(strings.Split(args, " ")) {
		splitArgs := strings.Split(args, " ")
		if len(splitArgs) == 0 {
			fmt.Println("[!] No time provided")
			zap.L().Fatal("no_time_provided",
				zap.String("message", "no time provided"),
			)
			os.Exit(1)
		} else if len(splitArgs) < 3 {
			fmt.Println("[!] Invalid time format")
			zap.L().Fatal("invalid_time_format",
				zap.String("message", "invalid time format"),
			)
			os.Exit(1)
		}

		// Handle 'last 24 hours' format
		var (
			tense    string
			amount   int
			duration time.Duration
		)
		for _, arg := range splitArgs {
			if isPasteTense(arg) {
				tense = arg
			} else if isNumber(arg) {
				amount, _ = strconv.Atoi(arg)
			} else if isDuration(arg) {
				duration = getDuration(arg)
			}
		}

		if tense == "" {
			fmt.Println("[!] Invalid time format: tense not found")
			zap.L().Fatal("invalid_time_format",
				zap.String("message", "invalid time format"),
			)
			os.Exit(1)
		} else if amount == 0 {
			fmt.Println("[!] Invalid time format: amount not found")
			zap.L().Fatal("invalid_time_format",
				zap.String("message", "invalid time format"),
			)
			os.Exit(1)
		} else if duration == 0 {
			fmt.Println("[!] Invalid time format: duration not found")
			zap.L().Fatal("invalid_time_format",
				zap.String("message", "invalid time format"),
			)
			os.Exit(1)
		}

		// Return the appropriate time
		if tense == "last" {
			return time.Now().Add(-time.Duration(amount) * duration)
		} else if tense == "ago" {
			return time.Now().Add(-time.Duration(amount) * duration)
		}
	}

	// Handle possible formats 'May 01, 2025', '05-01-2025', '05/01/2025', '05/01/25', '05-01-25'
	var (
		t     time.Time
		err   error
		found bool
	)
	possible := []string{"01-02-2006", "01/02/2006", "01/02/06", "01-02-06", "Jan 02, 2006", "Jan 2, 2006"}
	for _, format := range possible {
		t, err = time.Parse(format, args)
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("[!] Invalid time format")
		zap.L().Fatal("invalid_time_format",
			zap.String("message", "invalid time format"),
		)
		os.Exit(1)
	}

	// Convert UTC time to local time
	local, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("[!] Error loading local timezone")
		zap.L().Error("load_timezone",
			zap.String("message", "failed to load local timezone"),
			zap.Error(err),
		)
		return t
	}

	// Convert the parsed time to local time
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		local,
	)
}

func isPasteTense(value string) bool {
	for _, v := range []string{"last", "ago"} {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}

func isDuration(value string) bool {
	for _, v := range []string{"hour", "hours", "minute", "minutes", "second", "seconds", "day", "days", "week", "weeks", "month", "months", "year", "years"} {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}

func isNumber(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

func getDuration(timeBlock string) time.Duration {
	timeBlock = strings.TrimSpace(strings.ToLower(timeBlock))

	switch timeBlock {
	case "hour":
		return time.Hour
	case "hours":
		return time.Hour
	case "minute":
		return time.Minute
	case "minutes":
		return time.Minute
	case "second":
		return time.Second
	case "seconds":
		return time.Second
	case "day":
		return 24 * time.Hour
	case "days":
		return 24 * time.Hour
	case "week":
		return 7 * 24 * time.Hour
	case "weeks":
		return 7 * 24 * time.Hour
	case "month":
		return 30 * 24 * time.Hour
	case "months":
		return 30 * 24 * time.Hour
	case "year":
		return 365 * 24 * time.Hour
	case "years":
		return 365 * 24 * time.Hour
	default:
		fmt.Printf("[!] Unknown duration: %s", timeBlock)
		zap.L().Fatal("unknown_duration",
			zap.String("message", "unknown duration"),
			zap.String("duration", timeBlock),
		)
		os.Exit(1)
	}
	return 0
}

func containsMonth(arr []string) bool {
	for _, v := range arr {
		if isMonth(v) {
			return true
		}
	}
	return false
}

func isMonth(value string) bool {
	for _, v := range []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"} {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}
