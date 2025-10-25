package models

// LogEntry represents the incoming log format
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	IP        string `json:"ip"`
	URL       string `json:"url"`
}
