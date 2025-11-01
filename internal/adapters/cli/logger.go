package cli

import (
	"log"
	"sync"

	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// ANSI color codes for terminal output
const (
	colorRed    = "\033[91m"
	colorGreen  = "\033[92m"
	colorYellow = "\033[93m"
	colorReset  = "\033[0m"
)

// LoggerReporter is an output adapter for reporting progress and events.
type LoggerReporter struct {
	mu sync.Mutex
}

// NewLoggerReporter creates a new logger reporter.
func NewLoggerReporter() *LoggerReporter {
	return &LoggerReporter{}
}

// ReportProgress reports the current progress.
func (lr *LoggerReporter) ReportProgress(progress map[string]float64, isComplete bool) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	for name, pct := range progress {
		// Log all progress updates for real-time visibility of multiple concurrent conversions
		log.Printf("%s: %.1f%%", name, pct)
	}
}

// ReportConversionStart reports that a conversion has started.
func (lr *LoggerReporter) ReportConversionStart(filename string, hasSubs bool) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	subsInfo := "(no subtitles)"
	if hasSubs {
		subsInfo = "(with subtitles)"
	}
	log.Printf("%s[START]%s %s %s", colorGreen, colorReset, filename, subsInfo)
}

// ReportConversionFinish reports that a conversion has finished.
func (lr *LoggerReporter) ReportConversionFinish(filename string, outputPath string, success bool) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if success {
		log.Printf("%s[FINISH]%s %s -> %s", colorGreen, colorReset, filename, outputPath)
	} else {
		log.Printf("%s[FAILED]%s %s", colorRed, colorReset, filename)
	}
}

// ReportError reports an error with red color highlighting.
func (lr *LoggerReporter) ReportError(message string) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	log.Printf("%s[ERROR]%s %s", colorRed, colorReset, message)
}

var _ ports.ProgressReporterPort = (*LoggerReporter)(nil)
