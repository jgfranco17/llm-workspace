package tools

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentTime(t *testing.T) {
	before := time.Now().Truncate(time.Second)

	result, err := GetCurrentTime(context.Background(), Parameters{})

	after := time.Now().Add(time.Second)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	parsed, parseErr := time.Parse(time.RFC3339, result)
	assert.NoError(t, parseErr, "result should be a valid RFC3339 timestamp")
	assert.False(t, parsed.Before(before), "timestamp should not predate the call")
	assert.False(t, parsed.After(after), "timestamp should not postdate the call")
}

func TestGetCurrentTime_IgnoresParams(t *testing.T) {
	cases := []struct {
		name   string
		params Parameters
	}{
		{name: "nil params", params: nil},
		{name: "empty params", params: Parameters{}},
		{name: "irrelevant params", params: Parameters{"foo": "bar"}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := GetCurrentTime(context.Background(), testCase.params)
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
		})
	}
}
