package nginx

type LogEntry struct {
	Route      string
	StatusCode int
}

func (l LogEntry) HasValidStatusCode() bool {
	return l.StatusCode >= 200 && l.StatusCode <= 599
}

func (l LogEntry) Has20xStatusCode() bool {
	return l.StatusCode >= 200 && l.StatusCode <= 299
}

func (l LogEntry) Has30xStatusCode() bool {
	return l.StatusCode >= 300 && l.StatusCode <= 399
}

func (l LogEntry) Has40xStatusCode() bool {
	return l.StatusCode >= 400 && l.StatusCode <= 499
}

func (l LogEntry) Has50xStatusCode() bool {
	return l.StatusCode >= 500 && l.StatusCode <= 599
}

func (l LogEntry) HasErrorStatusCode() bool {
	return l.Has50xStatusCode()
}
