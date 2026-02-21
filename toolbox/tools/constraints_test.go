package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitOutput(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		limit int
		want  string
	}{
		{
			name:  "shorter than limit",
			text:  "hello",
			limit: 10,
			want:  "hello",
		},
		{
			name:  "exactly at limit",
			text:  "hello",
			limit: 5,
			want:  "hello",
		},
		{
			name:  "longer than limit",
			text:  "hello world",
			limit: 5,
			want:  "hello",
		},
		{
			name:  "empty string",
			text:  "",
			limit: 10,
			want:  "",
		},
		{
			name:  "limit of zero",
			text:  "hello",
			limit: 0,
			want:  "",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := limitOutput(testCase.text, testCase.limit)
			assert.Equal(t, testCase.want, got)
		})
	}
}
