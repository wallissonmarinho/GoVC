package ports

// ProgressReporterPort is the output port for reporting progress.
// (Right side of hexagon: core → UI/logging)
type ProgressReporterPort interface {
	ReportProgress(progress map[string]float64, isComplete bool)
	ReportConversionStart(filename string, hasSubs bool)
	ReportConversionFinish(filename string, outputPath string, success bool)
	ReportError(message string)
}
