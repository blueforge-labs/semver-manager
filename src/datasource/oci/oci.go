package oci

import (
	"github.com/blueforge-labs/semver-manager/src/models"
	"github.com/blueforge-labs/semver-manager/src/utils"
)

type OciCLient struct {
	config *utils.DatasourceConfig
}

func NewFetcher(config *utils.DatasourceConfig) *OciCLient {
	return &OciCLient{config: config}
}

func (g *OciCLient) FetchTags() ([]models.Version, error) {
	return []models.Version{}, nil
}
