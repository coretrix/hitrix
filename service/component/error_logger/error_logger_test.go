package errorlogger

import "testing"

func TestBuildErrorKeySource(t *testing.T) {
	tests := []struct {
		name      string
		errorCode string
		scopeID   string
		expected  string
	}{
		{
			name:     "default",
			expected: "web-api:/app/main.go:42",
		},
		{
			name:      "with error code",
			errorCode: "WORKSTATION_HEALTHCHECK_FAILED",
			expected:  "web-api:/app/main.go:42:WORKSTATION_HEALTHCHECK_FAILED",
		},
		{
			name:      "with error code and scope",
			errorCode: "WORKSTATION_HEALTHCHECK_FAILED",
			scopeID:   "machine-001",
			expected:  "web-api:/app/main.go:42:WORKSTATION_HEALTHCHECK_FAILED:machine-001",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := buildErrorKeySource("web-api", "/app/main.go", 42, test.errorCode, test.scopeID)
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
