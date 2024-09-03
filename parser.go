package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]interface{}
	Raw       string // Original log line for fallback
}

func parseLogrus(logLine string) (*LogEntry, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(logLine), &data); err != nil {
		return nil, err
	}

	level, _ := data["level"].(string)
	message, _ := data["msg"].(string)
	timestampStr, _ := data["time"].(string)

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Timestamp: timestamp,
		Level:     strings.ToUpper(level),
		Message:   message,
		Fields:    data,
		Raw:       logLine,
	}, nil
}

func parseZap(logLine string) (*LogEntry, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(logLine), &data); err != nil {
		return nil, err
	}

	level, _ := data["level"].(string)
	message, _ := data["msg"].(string)
	timestampStr, _ := data["ts"].(float64)

	timestamp := time.Unix(int64(timestampStr), 0)

	return &LogEntry{
		Timestamp: timestamp,
		Level:     strings.ToUpper(level),
		Message:   message,
		Fields:    data,
		Raw:       logLine,
	}, nil
}

func parseKlog(logLine string) (*LogEntry, error) {
	// Example klog format: "I0915 15:04:05.123456 12345 file.go:67] This is a klog message"
	// Level: I for Info, W for Warning, E for Error, F for Fatal

	if len(logLine) < 20 {
		return nil, fmt.Errorf("log line too short")
	}

	level := ""
	switch logLine[0] {
	case 'I':
		level = "INFO"
	case 'W':
		level = "WARNING"
	case 'E':
		level = "ERROR"
	case 'F':
		level = "FATAL"
	default:
		return nil, fmt.Errorf("unknown klog level")
	}

	timestampStr := logLine[1:16]
	timestamp, err := time.Parse("0102 15:04:05.000000", timestampStr)
	if err != nil {
		return nil, err
	}

	message := logLine[22:] // Skipping past the timestamp and file info

	return &LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
		Fields:    map[string]interface{}{"source": "klog"},
		Raw:       logLine,
	}, nil
}

func parseSlog(logLine string) (*LogEntry, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(logLine), &data); err != nil {
		return nil, err
	}

	level, _ := data["level"].(string)
	message, _ := data["msg"].(string)
	timestampStr, _ := data["time"].(string)

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Timestamp: timestamp,
		Level:     strings.ToUpper(level),
		Message:   message,
		Fields:    data,
		Raw:       logLine,
	}, nil
}
