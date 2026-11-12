package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/blueforge-labs/semver-manager/src/cmd/fetch"
	"github.com/blueforge-labs/semver-manager/src/cmd/filter"
	"github.com/blueforge-labs/semver-manager/src/cmd/increment"
)

func NewRootCommand(output io.Writer) *cli.Command {
	cmd := &cli.Command{
		Name:        "smgr",
		Usage:       "Manage Semantic Versioning compliant versions.",
		Description: `Manage Semantic Versioning compliant versions and integrate with popular repository and registry platform to facilitate the task.`,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "dry-run", Usage: "Execute the command in dry-run mode"},
			&cli.IntFlag{Name: "v", Usage: "Set logging verbosity"},
			&cli.StringFlag{Name: "prefix", Usage: "Prefix used for version tags and release output"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return cli.ShowRootCommandHelp(cmd)
		},
		Writer:    output,
		ErrWriter: os.Stderr,
	}

	filterArgs := &filter.FilterArgs{}
	filterCmd := filter.NewFilterCommand(filterArgs)
	fetchCmd := fetch.NewFetchCommand(filterArgs)
	cmd.Commands = []*cli.Command{filterCmd, fetchCmd, increment.NewIncrementCommand()}

	return cmd
}

func main() {
	rootCmd := NewRootCommand(os.Stdout)

	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
