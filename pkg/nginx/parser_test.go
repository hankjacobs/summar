package nginx

import "testing"

const (
	ValidLine                = `10.10.180.161 - 50.112.166.232, 192.33.28.238 - - - [02/Aug/2015:15:56:14 +0000]  https https https "GET /our-products HTTP/1.1" 200 35967 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36" `
	InvalidLine1             = `blah blah blah`
	InvalidRequestFieldsLine = `blah " blah " blah " blah " blah " blah " blah `
	InvalidStatusFieldsLine  = `blah " b la h " blah " blah " blah " blah " blah `
	InvalidStatusCodeLine    = `10.10.180.161 - 50.112.166.232, 192.33.28.238 - - - [02/Aug/2015:15:56:14 +0000]  https https https "GET /our-products HTTP/1.1" BAD 35967 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36" `
)

func TestParseValidLine(t *testing.T) {

	expected := LogEntry{Route: "/our-products", StatusCode: 200}
	got, err := ParseLogEntry(ValidLine)

	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if got != expected {
		t.Fatalf("Got %v but expected %v", got, expected)
	}
}

func TestParseInvalidLine(t *testing.T) {
	_, err := ParseLogEntry(InvalidLine1)

	if err != ErrInvalidLog {
		t.Fatalf("Got error %v but expected %v", err, ErrInvalidLog)
	}
}

func TestParseInvalidRequestFieldsLine(t *testing.T) {
	_, err := ParseLogEntry(InvalidRequestFieldsLine)

	if err != ErrInvalidLog {
		t.Fatalf("Got error %v but expected %v", err, ErrInvalidLog)
	}
}

func TestParseInvalidStatusFieldsLine(t *testing.T) {
	_, err := ParseLogEntry(InvalidStatusFieldsLine)

	if err != ErrInvalidLog {
		t.Fatalf("Got error %v but expected %v", err, ErrInvalidLog)
	}
}

func TestParseInvalidStatusCodeLine(t *testing.T) {
	_, err := ParseLogEntry(InvalidStatusCodeLine)

	if err == nil {
		t.Fatal("Got no error but expected error")
	}
}
