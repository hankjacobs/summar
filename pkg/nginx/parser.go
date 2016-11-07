package nginx

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrInvalidLog error returned for invalid log lines
	ErrInvalidLog = fmt.Errorf("invalid log line")
)

// ParseLogEntry parses a log entry string into a LogEntry
func ParseLogEntry(log string) (LogEntry, error) {

	substrs := strings.SplitN(log, `"`, -1)
	if len(substrs) != 7 {
		return LogEntry{}, ErrInvalidLog
	}

	reqFields := strings.Fields(substrs[1])
	if len(reqFields) != 3 {
		return LogEntry{}, ErrInvalidLog
	}

	statusFields := strings.Fields(substrs[2])
	if len(statusFields) != 2 {
		return LogEntry{}, ErrInvalidLog
	}

	route := reqFields[1]
	status, convErr := strconv.Atoi(statusFields[0])
	if convErr != nil {
		return LogEntry{}, convErr
	}

	return LogEntry{route, status}, nil
}
