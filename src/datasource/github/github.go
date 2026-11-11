package github

import (
	"github.com/bluepr-nt/semver-manager/src/models"
	"github.com/bluepr-nt/semver-manager/src/utils"
)

type GithubClient struct {
	config *utils.DatasourceConfig
}

func NewFetcher(config *utils.DatasourceConfig) *GithubClient {
	return &GithubClient{config: config}
}

func (g *GithubClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
