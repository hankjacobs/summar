package nginx

// LogEntry represents the values within a single access log file line
type LogEntry struct {
	Route      string
	StatusCode int
}

// HasValidStatusCode whether the log entry has a valid status code
func (l LogEntry) HasValidStatusCode() bool {
	return l.StatusCode >= 200 && l.StatusCode <= 599
}

// Has20xStatusCode whether the log entry has a 2xx status code
func (l LogEntry) Has20xStatusCode() bool {
	return l.StatusCode >= 200 && l.StatusCode <= 299
}

// Has30xStatusCode whether the log entry has a 3xx status code
func (l LogEntry) Has30xStatusCode() bool {
	return l.StatusCode >= 300 && l.StatusCode <= 399
}

// Has40xStatusCode whether the log entry has a 4xx status code
func (l LogEntry) Has40xStatusCode() bool {
	return l.StatusCode >= 400 && l.StatusCode <= 499
}

// Has50xStatusCode whether the log entry has a 5xx status code
func (l LogEntry) Has50xStatusCode() bool {
	return l.StatusCode >= 500 && l.StatusCode <= 599
}
