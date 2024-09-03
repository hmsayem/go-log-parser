package main

import (
	"fmt"
)

func main() {
	// Sample logs (in practice, read from a file or other source)
	logLines := []string{
		`{"level":"info","msg":"This is a logrus log","time":"2023-09-01T00:00:00Z"}`,
		`{"level":"error","msg":"This is a zap log","ts":1693574400}`,
		`I0915 15:04:05.123456 12345 file.go:67] This is a klog message`,
		`{"level":"warn","msg":"This is a slog log","time":"2023-09-01T00:00:00Z"}`,
	}

	filtered := filterLogs(logLines, "error", []string{"zap", "klog"})
	for _, log := range filtered {
		fmt.Println(log)
	}
}
