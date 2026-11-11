package gitlab

import (
	"github.com/bluepr-nt/semver-manager/src/models"
	"github.com/bluepr-nt/semver-manager/src/utils"
)

type GitlabClient struct {
	config *utils.DatasourceConfig
}

func NewFetcher(config *utils.DatasourceConfig) *GitlabClient {
	return &GitlabClient{config: config}
}

func (g *GitlabClient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
