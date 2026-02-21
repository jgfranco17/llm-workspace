package tools

import (
	"context"
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

func TestGetRequiredParam(t *testing.T) {
	cases := []struct {
		name         string
		params       Parameters
		key          string
		wantErr      bool
		wantErrMatch string
		wantValue    string
	}{
		{
			name:      "key present with value",
			params:    Parameters{"cmd": "echo hi"},
			key:       "cmd",
			wantValue: "echo hi",
		},
		{
			name:         "key missing",
			params:       Parameters{},
			key:          "cmd",
			wantErr:      true,
			wantErrMatch: "missing required parameter: cmd",
		},
		{
			name:         "key present but empty string",
			params:       Parameters{"cmd": ""},
			key:          "cmd",
			wantErr:      true,
			wantErrMatch: "missing required parameter: cmd",
		},
		{
			name:         "key present but only whitespace",
			params:       Parameters{"cmd": "   "},
			key:          "cmd",
			wantErr:      true,
			wantErrMatch: "missing required parameter: cmd",
		},
		{
			name:      "value preserved as-is including surrounding whitespace",
			params:    Parameters{"cmd": " echo "},
			key:       "cmd",
			wantValue: " echo ",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := getRequiredParam(context.Background(), testCase.params, testCase.key)
			if testCase.wantErr {
				assert.ErrorContains(t, err, testCase.wantErrMatch)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantValue, got)
			}
		})
	}
}
