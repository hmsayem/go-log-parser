package main

import (
	"encoding/json"
	"regexp"
	"strings"
)

// LogType represents the type of logging library.
type LogType string

const (
	LogTypeKlog    LogType = "klog"
	LogTypeSlog    LogType = "slog"
	LogTypeLogrus  LogType = "logrus"
	LogTypeZap     LogType = "zap"
	LogTypeUnknown LogType = "unknown"
)

// IdentifyLogType tries to identify the log format based on the log line.
func logType(logLine string) LogType {
	// Check for klog format
	klogRegex := regexp.MustCompile(`^[IWEF]\d{4} \d{2}:\d{2}:\d{2}\.\d{6}`)
	if klogRegex.MatchString(logLine) {
		return LogTypeKlog
	}

	// Check for JSON format (zap, slog)
	if strings.HasPrefix(logLine, "{") && strings.HasSuffix(logLine, "}") {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(logLine), &jsonMap); err == nil {
			// Check for zap specific fields
			if _, ok := jsonMap["ts"]; ok {
				if _, ok := jsonMap["caller"]; ok {
					return LogTypeZap
				}
			}
			// Check for slog specific fields
			if _, ok := jsonMap["time"]; ok {
				if _, ok := jsonMap["level"]; ok {
					return LogTypeSlog
				}
			}
		}
	}

	logrusRegex := regexp.MustCompile(`time=".*" level=.* msg=".*"`)
	if logrusRegex.MatchString(logLine) {
		return LogTypeLogrus
	}

	return LogTypeUnknown
}

func filterLogs(logLines []string, logLevels []string, keywords []string) []string {
	var filteredLogs []string

	for _, line := range logLines {
		var entry *LogEntry
		var err error

		// Try parsing with each known log format
		if logType(line) == LogTypeLogrus {
			entry, err = parseLogrus(line)
		} else if logType(line) == LogTypeZap {
			entry, err = parseZap(line)
		} else if logType(line) == LogTypeKlog {
			entry, err = parseKlog(line)
		} else if logType(line) == LogTypeSlog {
			entry, err = parseSlog(line)
		}

		// Add more parsers for other log libraries here
		if err != nil || entry == nil {
			continue
		}

		// Apply log level and keyword filters
		if logLevelsMatched(entry, logLevels) && containsKeywords(entry, keywords) {
			filteredLogs = append(filteredLogs, line)
		}
	}

	return filteredLogs
}

func logLevelsMatched(entry *LogEntry, levels []string) bool {
	if levels == nil {
		return true
	}

	for i, _ := range levels {
		if levels[i] == entry.Level {
			return true
		}
	}
	return false
}

func containsKeywords(entry *LogEntry, keywords []string) bool {
	if keywords == nil {
		return true
	}

	for _, keyword := range keywords {
		if strings.Contains(entry.Message, keyword) {
			return true
		}
	}
	return false
}
