package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.params.Get(tc.key)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.wantErrMatch)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, got)
			}
		})
	}
}
