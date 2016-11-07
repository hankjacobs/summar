package nginx

import "testing"

func TestHasValidStatusCode(t *testing.T) {
	for i := 200; i <= 599; i++ {
		entry := LogEntry{"/test", i}
		if !entry.HasValidStatusCode() {
			t.Errorf("Entry Status Code (%v) was invalid when it should have been valid", entry.StatusCode)
		}
	}

	for i := 100; i <= 199; i++ {
		entry := LogEntry{"/test", i}
		if entry.HasValidStatusCode() {
			t.Errorf("Entry Status Code (%v) was valid when it should not have been valid", entry.StatusCode)
		}
	}
}

func TestHas20xStatusCode(t *testing.T) {
	for i := 200; i <= 299; i++ {
		entry := LogEntry{"/test", i}
		if !entry.Has20xStatusCode() {
			t.Errorf("Entry Status Code (%v) not 20x when it should be", entry.StatusCode)
		}
	}

	for i := 100; i <= 199; i++ {
		entry := LogEntry{"/test", i}
		if entry.Has20xStatusCode() {
			t.Errorf("Entry Status Code (%v) is 20x when it should not be", entry.StatusCode)
		}
	}
}

func TestHas30xStatusCode(t *testing.T) {
	for i := 300; i <= 399; i++ {
		entry := LogEntry{"/test", i}
		if !entry.Has30xStatusCode() {
			t.Errorf("Entry Status Code (%v) not 30x when it should be", entry.StatusCode)
		}
	}

	for i := 100; i <= 199; i++ {
		entry := LogEntry{"/test", i}
		if entry.Has30xStatusCode() {
			t.Errorf("Entry Status Code (%v) is 30x when it should not be", entry.StatusCode)
		}
	}
}

func TestHas40xStatusCode(t *testing.T) {
	for i := 400; i <= 499; i++ {
		entry := LogEntry{"/test", i}
		if !entry.Has40xStatusCode() {
			t.Errorf("Entry Status Code (%v) not 40x when it should be", entry.StatusCode)
		}
	}

	for i := 100; i <= 199; i++ {
		entry := LogEntry{"/test", i}
		if entry.Has40xStatusCode() {
			t.Errorf("Entry Status Code (%v) is 40x when it should not be", entry.StatusCode)
		}
	}
}

func TestHas50xStatusCode(t *testing.T) {
	for i := 500; i <= 599; i++ {
		entry := LogEntry{"/test", i}
		if !entry.Has50xStatusCode() {
			t.Errorf("Entry Status Code (%v) not 50x when it should be", entry.StatusCode)
		}
	}

	for i := 100; i <= 199; i++ {
		entry := LogEntry{"/test", i}
		if entry.Has50xStatusCode() {
			t.Errorf("Entry Status Code (%v) is 50x when it should not be", entry.StatusCode)
		}
	}
}
