package domain

// ConversionResult represents the result of a video conversion.
type ConversionResult struct {
	Video   *Video
	Success bool
	Error   string
	Message string
}

// NewSuccessResult creates a successful conversion result.
func NewSuccessResult(video *Video) *ConversionResult {
	return &ConversionResult{
		Video:   video,
		Success: true,
		Message: "Conversion completed successfully",
	}
}

// NewErrorResult creates a failed conversion result.
func NewErrorResult(video *Video, err string) *ConversionResult {
	return &ConversionResult{
		Video:   video,
		Success: false,
		Error:   err,
		Message: "Conversion failed",
	}
}
