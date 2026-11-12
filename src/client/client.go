package client

import (
	"context"
	"fmt"

	datasource "github.com/blueforge-labs/semver-manager/src/datasource/utils"
	"github.com/blueforge-labs/semver-manager/src/models"
	"github.com/blueforge-labs/semver-manager/src/pkg/filter"
	"github.com/blueforge-labs/semver-manager/src/pkg/increment"
	"github.com/blueforge-labs/semver-manager/src/pkg/push"
	"github.com/blueforge-labs/semver-manager/src/pkg/tagprefix"
	"github.com/blueforge-labs/semver-manager/src/pkg/validate"
)

// Config contains repository settings and the tag namespace used by the client.
type Config struct {
	Owner      string
	Repository string
	Token      string
	Platform   string
	Prefix     string
	LocalPath  string
}

// Client exposes the semver-manager operations to Go applications.
type Client struct {
	config Config
}

// New creates a client. GitHub is used when Platform is empty.
func New(config Config) *Client {
	if config.Platform == "" {
		config.Platform = "github"
	}
	return &Client{config: config}
}

// Fetch reads repository tags, removes the configured prefix, and parses SemVer tags.
func (c *Client) Fetch() ([]models.Version, error) {
	datasource := datasource.NewSemverSvc(c.config.Platform, c.config.Token)
	tags, err := datasource.FetchSemverTags(c.config.Owner, c.config.Repository)
	if err != nil {
		return nil, err
	}
	return c.parseTags(tags), nil
}

func (c *Client) parseTags(tags []string) []models.Version {
	return tagprefix.Parse(tags, c.config.Prefix)
}

// Filter applies the supplied filters to versions.
func (c *Client) Filter(versions []models.Version, filters ...filter.FilterFunc) (models.VersionSlice, error) {
	return filter.ApplyFilters(versions, filters...)
}

// Increment calculates the next version in the requested stream.
func (c *Client) Increment(sourceVersions []models.Version, stream models.VersionPattern, level models.Increment) (models.Version, error) {
	return increment.IncrementVersion(sourceVersions, stream, level)
}

// FormatRelease prefixes a semantic version for use as a repository tag name.
func (c *Client) FormatRelease(version models.Version) string {
	return tagprefix.Format(version, c.config.Prefix)
}

// Push creates a Git tag targeting a commit SHA, a GitHub ref, or local HEAD when target is empty.
func (c *Client) Push(ctx context.Context, version models.Version, target string) error {
	if c.config.Platform != "github" {
		return fmt.Errorf("tag push is not supported for platform %q", c.config.Platform)
	}
	return push.PushTag(ctx, push.Config{
		Owner:      c.config.Owner,
		Repository: c.config.Repository,
		Token:      c.config.Token,
		LocalPath:  c.config.LocalPath,
	}, c.FormatRelease(version), target)
}

// Validate validates a version, stripping the configured prefix first.
func (c *Client) Validate(version, validatorType string) (bool, error) {
	if c.config.Prefix != "" {
		stripped, matches := tagprefix.Strip(version, c.config.Prefix)
		if !matches {
			return false, fmt.Errorf("version %q does not have configured prefix %q", version, c.config.Prefix)
		}
		version = stripped
	}
	validator := validate.NewSemverValidator(validatorType)
	if validator == nil {
		return false, fmt.Errorf("unsupported validator type %q", validatorType)
	}
	return validate.IsSemverValid(version, validator)
}
