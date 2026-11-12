package filter

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestNewFilterCommand(t *testing.T) {
	tests := []struct {
		name          string
		inputArgs     []string
		expectedError bool
		expectedOut   string
	}{
		{
			name:        "Provided single version highest",
			inputArgs:   []string{"--versions", "1.2.3", "--highest"},
			expectedOut: "1.2.3",
		},
		{
			name:        "Provided multiple versions highest",
			inputArgs:   []string{"--versions", "1.2.3, 1.1.1", "--highest"},
			expectedOut: "1.2.3",
		},
		{
			name:        "Provided multiple versions highest with some bad versions",
			inputArgs:   []string{"--versions", "1.2.3, 1.1.1, bad.version", "--highest"},
			expectedOut: "1.2.3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filterArgs := &FilterArgs{}
			filtercmd := NewFilterCommand(filterArgs)
			output, err := executeCommand(filtercmd, test.inputArgs...)

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedOut, output)
		})
	}
}

func executeCommand(cmd *cli.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.Writer = buf
	cmd.ErrWriter = buf
	arguments := append([]string{cmd.Name}, args...)
	err := cmd.Run(context.Background(), arguments)
	return strings.TrimSpace(buf.String()), err
}
