package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestNewRootCommand(t *testing.T) {

	t.Run("Base Tests", func(t *testing.T) {
		output := &bytes.Buffer{}
		cmd := NewRootCommand(output)
		assert.NotEmpty(t, cmd, "NewRootCommand() should not return nil command")
		assert.Equal(t, "smgr", cmd.Name, "NewRootCommand() should return command with name 'smgr'")
		assert.True(t, hasFlag(cmd.Flags, "dry-run"), "NewRootCommand() should have 'dry-run' flag")
		assert.True(t, hasFlag(cmd.Flags, "prefix"), "NewRootCommand() should have 'prefix' flag")
	})
}

func TestPrefixFlagAppliesToSubcommands(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "filter",
			args: []string{"smgr", "filter", "--prefix", "release/", "--versions", "release/1.2.3,release/1.1.1", "--highest"},
			want: "release/1.2.3\n",
		},
		{
			name: "increment",
			args: []string{"smgr", "--prefix", "v", "increment", "--source-versions", "v1.2.3", "--level", "patch"},
			want: "v1.2.4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			cmd := NewRootCommand(output)
			require.NoError(t, cmd.Run(context.Background(), test.args))
			assert.Equal(t, test.want, output.String())
		})
	}
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
