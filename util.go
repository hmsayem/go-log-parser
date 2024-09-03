package main

import "strings"

func filterLogs(logLines []string, logLevel string, keywords []string) []string {
	var filteredLogs []string

	for _, line := range logLines {
		var entry *LogEntry
		var err error

		// Try parsing with each known log format
		if strings.Contains(line, "logrus") {
			entry, err = parseLogrus(line)
		} else if strings.Contains(line, "zap") {
			entry, err = parseZap(line)
		} else if strings.Contains(line, "klog") || strings.HasPrefix(line, "I") || strings.HasPrefix(line, "W") || strings.HasPrefix(line, "E") || strings.HasPrefix(line, "F") {
			entry, err = parseKlog(line)
		} else if strings.Contains(line, "slog") {
			entry, err = parseSlog(line)
		}

		// Add more parsers for other log libraries here

		if err != nil || entry == nil {
			continue
		}

		// Apply log level and keyword filters
		if entry.Level == strings.ToUpper(logLevel) && containsKeywords(entry, keywords) {
			filteredLogs = append(filteredLogs, line)
		}
	}

	return filteredLogs
}

func containsKeywords(entry *LogEntry, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(entry.Message, keyword) {
			return true
		}
	}
	return false
}
