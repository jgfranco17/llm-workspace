package logging

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewJSONLogger_Level(t *testing.T) {
	cases := []struct {
		name      string
		level     string
		wantLevel logrus.Level
	}{
		{
			name:      "trace level",
			level:     "trace",
			wantLevel: logrus.TraceLevel,
		},
		{
			name:      "debug level",
			level:     "debug",
			wantLevel: logrus.DebugLevel,
		},
		{
			name:      "info level",
			level:     "info",
			wantLevel: logrus.InfoLevel,
		},
		{
			name:      "warn level",
			level:     "warn",
			wantLevel: logrus.WarnLevel,
		},
		{
			name:      "error level",
			level:     "error",
			wantLevel: logrus.ErrorLevel,
		},
		{
			name:      "uppercase INFO is accepted",
			level:     "INFO",
			wantLevel: logrus.InfoLevel,
		},
		{
			name:      "invalid level defaults to info",
			level:     "invalid-level",
			wantLevel: logrus.InfoLevel,
		},
		{
			name:      "empty level defaults to info",
			level:     "",
			wantLevel: logrus.InfoLevel,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			logger := NewJSONLogger(&bytes.Buffer{}, testCase.level)
			assert.Equal(t, testCase.wantLevel, logger.GetLevel())
		})
	}
}

func TestNewJSONLogger_WritesToProvidedWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, "info")

	logger.Info("test message")

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "test message")
}

func TestNewJSONLogger_OutputIsJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, "info")

	logger.Info("json output test")

	var parsed map[string]any
	err := json.Unmarshal(buf.Bytes(), &parsed)
	assert.NoError(t, err, "logger output should be valid JSON")
}

func TestNewJSONLogger_InvalidLevel_LogsWarning(t *testing.T) {
	var buf bytes.Buffer
	NewJSONLogger(&buf, "not-a-level")

	output := buf.String()
	assert.Contains(t, output, "not-a-level")
}
