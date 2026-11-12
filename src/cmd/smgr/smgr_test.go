package main

import (
	"bytes"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestNewRootCommand(t *testing.T) {

	t.Run("Base Tests", func(t *testing.T) {
		output := &bytes.Buffer{}
		cmd := NewRootCommand(output)
		assert.NotEmpty(t, cmd, "NewRootCommand() should not return nil command")
		assert.Equal(t, "smgr", cmd.Name, "NewRootCommand() should return command with name 'smgr'")
		assert.True(t, hasFlag(cmd.Flags, "dry-run"), "NewRootCommand() should have 'dry-run' flag")
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
