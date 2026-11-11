package fetch

import (
	"errors"

	"github.com/bluepr-nt/semver-manager/src/datasource/github"
	"github.com/bluepr-nt/semver-manager/src/datasource/gitlab"
	"github.com/bluepr-nt/semver-manager/src/datasource/oci"
	"github.com/bluepr-nt/semver-manager/src/models"
	"github.com/bluepr-nt/semver-manager/src/utils"
)

type Fetcher interface {
	FetchTags() ([]models.Version, error)
}

func NewFetcher(config *utils.DatasourceConfig) (Fetcher, error) {
	switch config.Platform {
	case "github":
		return github.NewFetcher(config), nil
	case "gitlab":
		return gitlab.NewFetcher(config), nil
	case "oci":
		return oci.NewFetcher(config), nil
	default:
		return nil, errors.New("unsupported platform")
	}
}
