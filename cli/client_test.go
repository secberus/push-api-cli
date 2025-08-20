package cli

import (
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {

}

func TestDecodeYAML(t *testing.T) {

}

func TestDecodeCSV(t *testing.T) {

}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		expected FileFormat
	}{
		{
			name:     "Valid JSON",
			content:  `{"name": "Alice", "age": 30}`,
			filename: "data.json",
			expected: FormatJSON,
		},
		{
			name:     "Valid YAML",
			content:  "name: Bob\nage: 40",
			filename: "config.yaml",
			expected: FormatYAML,
		},
		{
			name:     "Valid CSV",
			content:  "name,age\nCharlie,22",
			filename: "people.csv",
			expected: FormatCSV,
		},
		{
			name:     "Invalid content with .json extension",
			content:  "gibberish!",
			filename: "broken.json",
			expected: FormatJSON, // fallback to extension
		},
		{
			name:     "Invalid content with .yml extension",
			content:  "%%%notyaml",
			filename: "bad.yml",
			expected: FormatYAML, // fallback to extension
		},
		{
			name:     "Unknown content and extension",
			content:  "!!!???",
			filename: "unknown.bin",
			expected: FormatUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, _ := detectFormat(strings.NewReader(tt.content), tt.filename)
			if format != tt.expected {
				t.Errorf("expected format %s, got %s", tt.expected, format)
			}
		})
	}
}
