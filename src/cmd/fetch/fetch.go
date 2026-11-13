package fetch

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/blueforge-labs/semver-manager/src/cmd/filter"
	datasourceUtils "github.com/blueforge-labs/semver-manager/src/datasource/utils"
	"github.com/blueforge-labs/semver-manager/src/pkg/tagprefix"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type config struct {
	Token      string `san:"trim"`
	Repository string `san:"trim"`
	Owner      string `san:"trim"`
	Platform   string `san:"trim"`
	dryRun     bool
}

func NewFetchCommand(filterArgs *filter.FilterArgs) *cli.Command {
	if filterArgs == nil {
		filterArgs = &filter.FilterArgs{}
	}
	config := &config{}
	fetchCmd := &cli.Command{
		Name:  "fetch",
		Usage: "Fetch semver tags from a repository.",
		Description: `Fetch semver tags from a repository and
sorted by highest version first. Fetch also supports all
the filters from the filter command. If the --versions
flag is set, the versions passed will be merged with the
fetched versions.`,
		Flags: filter.FilterFlags(filterArgs),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.dryRun = cmd.Bool("dry-run")
			filterArgs.Prefix = cmd.String("prefix")
			return RunFetchSemverTags(config, cmd, filterArgs)
		},
	}
	fetchCmd.Flags = append(fetchCmd.Flags,
		&cli.StringFlag{Name: "owner", Aliases: []string{"o"}, Usage: "The owner of the registry or repository", Destination: &config.Owner},
		&cli.StringFlag{Name: "repo", Aliases: []string{"r"}, Usage: "The repository or registry to fetch the Semver tags from", Destination: &config.Repository},
		&cli.StringFlag{Name: "token", Aliases: []string{"t"}, Usage: "The token to access the repository", Destination: &config.Token},
		&cli.StringFlag{Name: "platform", Aliases: []string{"p"}, Value: "github", Usage: "The platform to fetch the Semver from, options: github", Destination: &config.Platform},
	)

	return fetchCmd
}

func RunFetchSemverTags(config *config, cmd *cli.Command, filterArgs *filter.FilterArgs) error {

	datasource := newDatasource(config.dryRun, config.Platform, config.Token)

	logger := log.Logger.Output(os.Stderr).Level(zerolog.InfoLevel)
	if cmd.Int("v") >= 1 {
		logger = logger.Level(zerolog.DebugLevel)
	}
	logger.Debug().Msg("Fetching tags...")
	semverTags, err := datasource.FetchSemverTags(config.Owner, config.Repository)
	if err != nil {
		return err
	}
	logger.Debug().Int("count", len(semverTags)).Msgf("Fetched %d tags", len(semverTags))
	filterArgs.Versions = strings.Join(semverTags, " ")
	filteredTags, err := filter.Filter(filterArgs)
	if err != nil {
		return err
	}
	sort.SliceStable(filteredTags, func(i, j int) bool {
		if filteredTags[i].IsHigherThan(filteredTags[j]) {
			return true
		}
		if filteredTags[j].IsHigherThan(filteredTags[i]) {
			return false
		}
		return filteredTags[i].String() < filteredTags[j].String()
	})
	_, err = fmt.Fprintln(cmd.Writer, tagprefix.FormatVersions(filteredTags, filterArgs.Prefix))

	return err
}

func newDatasource(dryRun bool, platform, token string) datasourceUtils.Datasource {
	if len(platform) == 0 {
		platform = "github"
	}
	if dryRun {
		platform = "dry-run"
	}
	semverSvc := datasourceUtils.NewSemverSvc(platform, token)

	return semverSvc
}
