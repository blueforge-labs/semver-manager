package filter

import (
	"context"
	"fmt"

	"github.com/blueforge-labs/semver-manager/src/models"
	"github.com/blueforge-labs/semver-manager/src/pkg/filter"

	"github.com/urfave/cli/v3"
)

func NewFilterCommand(filterArgs *FilterArgs) *cli.Command {
	if filterArgs == nil {
		filterArgs = &FilterArgs{}
	}
	return &cli.Command{
		Name:        "filter",
		Usage:       "Filter is a CLI tool for filtering versions",
		Description: `Filter is a CLI tool for filtering versions using various criteria.`,
		Flags:       FilterFlags(filterArgs),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			semverTags, err := Filter(filterArgs)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.Writer, semverTags.String())
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func FilterFlags(filterArgs *FilterArgs) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "versions", Aliases: []string{"V"}, Usage: "Version list to filter", Destination: &filterArgs.Versions},
		&cli.StringFlag{Name: "stream", Aliases: []string{"s"}, Usage: "Filter by major, minor, patch, prerelease version and build metadata streams", Destination: &filterArgs.StreamFilter},
		&cli.BoolFlag{Name: "highest", Aliases: []string{"H"}, Usage: "Filter by highest version", Destination: &filterArgs.Highest},
	}
}

func Filter(filterArgs *FilterArgs) (models.VersionSlice, error) {
	versions := filter.GetValidVersions(filterArgs.Versions)

	filters := []filter.FilterFunc{}
	if filterArgs.StreamFilter != "" {
		pattern, err := models.ParseVersionPattern(filterArgs.StreamFilter)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter.VersionPatternFilter(pattern))
	}

	if filterArgs.Highest {
		filters = append(filters, filter.Highest())
	}

	semverTags, err := filter.ApplyFilters(versions, filters...)
	if err != nil {
		return nil, err
	}
	return semverTags, nil
}
