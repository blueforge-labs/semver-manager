package increment

import (
	"context"
	"fmt"
	"strings"

	"github.com/blueforge-labs/semver-manager/src/models"
	"github.com/blueforge-labs/semver-manager/src/pkg/increment"
	"github.com/blueforge-labs/semver-manager/src/pkg/tagprefix"

	"github.com/urfave/cli/v3"
)

type config struct {
	dryRun         bool
	incrementType  string
	sourceVersions string
	repository     string
	targetStream   string
	prefix         string
}

func NewIncrementCommand() *cli.Command {
	config := &config{}
	incrementCmd := &cli.Command{
		Name:  "increment",
		Usage: "Increment a version",
		Description: `
Increment a version according to one of the required flags --level or --target-stream, 
and any combination of the optional flags:

- Use --level to specify the increment level (major, minor, patch).
- Define the source with --repository, --source-stream, --source-version, or --source-versions. (Only --source-versions is currently implemented)

Increment a version according to the provided:
  - Increment level (major, minor, patch)
  - The source, any of: repository, source-stream, source-version, source-versions
  - The target stream, if specified, e.g. 1.2.* (optional)
  `,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "level", Aliases: []string{"l"}, Value: string(models.Patch), Usage: "The level of increment to perform: major, minor, patch", Destination: &config.incrementType},
			&cli.StringFlag{Name: "target-stream", Aliases: []string{"t"}, Usage: "The target stream to increment to e.g. 1.2.* (optional)", Destination: &config.targetStream},
			&cli.StringFlag{Name: "source-versions", Aliases: []string{"s"}, Usage: "The source versions to increment from e.g. \"0.0.0,1.0.0,1.1.0\" (optional)", Destination: &config.sourceVersions},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.dryRun = cmd.Bool("dry-run")
			config.prefix = cmd.String("prefix")
			level := cmd.String("level")
			targetStream := cmd.String("target-stream")

			if level == "" && targetStream == "" {
				config.incrementType = string(models.Patch)
			}

			return RunIncrement(config, cmd)
		},
	}

	return incrementCmd
}

func RunIncrement(config *config, cmd *cli.Command) error {
	sourceVersions := []models.Version{{}}
	if strings.TrimSpace(config.sourceVersions) != "" {
		sourceVersions = tagprefix.Parse([]string{config.sourceVersions}, config.prefix)
	}

	var err error
	var targetStream models.VersionPattern
	if config.targetStream != "" {
		targetStream, err = models.ParseVersionPattern(config.targetStream)
		if err != nil {
			return err
		}
	}

	newVersion, err := increment.IncrementVersion(sourceVersions, targetStream, models.Increment(config.incrementType))
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(cmd.Writer, tagprefix.Format(newVersion, config.prefix))
	return err
}
