package increment

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

type testFlag struct {
	name  string
	value string
}

func TestNewIncrementCommand(t *testing.T) {
	t.Run("Returns non-nil command", func(t *testing.T) {
		cmd := NewIncrementCommand()
		assert.NotNil(t, cmd)
	})

	t.Run("Command name is 'increment'", func(t *testing.T) {
		cmd := NewIncrementCommand()
		assert.Equal(t, "increment", cmd.Name)
	})

	t.Run("Command has expected flags", func(t *testing.T) {
		cmd := NewIncrementCommand()
		expectedFlags := []string{"level", "source-versions", "target-stream"}
		for _, expectedFlag := range expectedFlags {
			assert.True(t, hasFlag(cmd.Flags, expectedFlag))
		}
	})
}

func hasFlag(flags []cli.Flag, name string) bool {
	for _, flag := range flags {
		for _, flagName := range flag.Names() {
			if flagName == name {
				return true
			}
		}
	}
	return false
}

func TestNewIncrementCommandFlags(t *testing.T) {
	tests := []struct {
		name               string
		flags              []testFlag
		expectedNewVersion string
		expectedError      error
	}{
		{
			name: "Increment major version",
			flags: []testFlag{
				{name: "level", value: "major"},
			},
			expectedNewVersion: "1.0.0",
			expectedError:      nil,
		}, {
			name: "Increment minor version",
			flags: []testFlag{
				{name: "level", value: "minor"},
			},
			expectedNewVersion: "0.1.0",
			expectedError:      nil,
		},
		{
			name: "Increment patch version",
			flags: []testFlag{
				{name: "level", value: "patch"},
			},
			expectedNewVersion: "0.0.1",
			expectedError:      nil,
		},
		{
			name:               "Increment with no flags, default patch",
			flags:              []testFlag{},
			expectedNewVersion: "0.0.1",
			expectedError:      nil,
		},
		{
			name: "Increment major version with source versions",
			flags: []testFlag{
				{name: "level", value: "major"},
				{name: "source-versions", value: "0.0.0,1.0.0,0.1.0"},
			},
			expectedNewVersion: "2.0.0",
			expectedError:      nil,
		},
		{
			name: "Increment minor version with source versions and Prerelease target stream",
			flags: []testFlag{
				{name: "level", value: "minor"},
				{name: "source-versions", value: "0.0.0,1.0.0,0.1.0"},
				{name: "target-stream", value: "*.*.*-alpha.*"},
			},
			expectedNewVersion: "1.1.0-alpha.0",
			expectedError:      nil,
		},
		{
			name: "Increment minor version with source versions and Prerelease target stream",
			flags: []testFlag{
				{name: "level", value: "minor"},
				{name: "source-versions", value: "0.0.0,1.0.0-beta,0.1.0"},
				{name: "target-stream", value: "*.*.*-alpha.*"},
			},
			expectedNewVersion: "0.2.0-alpha.0",
			expectedError:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := new(bytes.Buffer)

			cmd := NewIncrementCommand()
			cmd.Writer = output
			cmd.ErrWriter = output
			args := []string{cmd.Name}

			for _, flag := range tt.flags {
				args = append(args, "--"+flag.name, flag.value)
			}

			err := cmd.Run(context.Background(), args)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedNewVersion, output.String())
		})
	}
}
