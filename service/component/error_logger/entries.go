package errorlogger

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Entry is the dev-panel representation of a stored log message.
type Entry struct {
	File       string
	Line       int
	AppName    string
	Service    string
	Request    string
	Message    string
	Stack      string
	Counter    int
	Time       string
	TicketLink string
	Timestamp  int64 `json:"-"`
}

// ParseEntries reconstructs message fields and occurrence metadata from a Redis hash.
func ParseEntries(data map[string]string) (map[string]*Entry, error) {
	entries := make(map[string]*Entry)

	for key, value := range data {
		if value == "" {
			continue
		}

		entryID, field := key, ""
		if lastSplitIndex := strings.LastIndex(key, ":"); lastSplitIndex > 0 {
			entryID, field = key[:lastSplitIndex], key[lastSplitIndex+1:]
		}

		if entries[entryID] == nil {
			entries[entryID] = &Entry{}
		}

		entry := entries[entryID]

		switch field {
		case "":
			var message ErrorMessage
			if err := json.Unmarshal([]byte(value), &message); err != nil {
				return nil, err
			}

			entry.Request, entry.Stack = string(message.Request), string(message.Stack)
			entry.File, entry.Line = message.File, message.Line
			entry.Message, entry.AppName, entry.Service = message.Message, message.AppName, message.Service
		case "time":
			entry.Timestamp, _ = strconv.ParseInt(strings.TrimSpace(value), 10, 64)
			entry.Time = time.Unix(entry.Timestamp, 0).String()
		case "counter":
			entry.Counter, _ = strconv.Atoi(strings.TrimSpace(value))
		case "ticket_link":
			entry.TicketLink = value
		}
	}

	return entries, nil
}
