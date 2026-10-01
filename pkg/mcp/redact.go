package mcp

import "regexp"

var sensitiveField = regexp.MustCompile(
	`(?im)(["']?(?:password|token|authorization|refresh[_-]?token|access[_-]?token|secret|api[_-]?key|cookie|set-cookie)` +
		`["']?\s*[:=]\s*)("(?:\\.|[^"\\])*"|'[^']*'|[^\r\n,}&]*)`,
)

var bearerToken = regexp.MustCompile(`(?i)(Bearer\s+)[a-z0-9._~+/=-]+`)

func maskSensitiveText(value string) string {
	return bearerToken.ReplaceAllString(sensitiveField.ReplaceAllString(value, "${1}***"), "${1}***")
}
