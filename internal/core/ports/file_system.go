package ports

// FileSystemPort is the output port for file system operations.
// (Right side of hexagon: core → file system)
type FileSystemPort interface {
	FileExists(path string) bool
	IsValidOutput(path string) bool
	RemoveFile(path string) error
	WriteLog(logPath string, lines []string) error
}
