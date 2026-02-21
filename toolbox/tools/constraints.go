package tools

const (
	// MaxOutputSize is the maximum number of bytes returned by a tool.
	MaxOutputSize = 1000

	// MaxFileSize is the maximum file size in bytes that ReadFile will accept.
	MaxFileSize = 100_000
)

func limitOutput(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	return text[:limit]
}
