package testutils

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/bluepr-nt/semver-manager/src/models"
)

func GetGithubToken() string {
	_ = godotenv.Load("fetch_test.yaml")
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("GITHUB_TOKEN is required for these tests. Authenticate with `gh auth login`, then set it with:\n" +
			"  fish: set -x GITHUB_TOKEN (gh auth token)\n" +
			"  bash/zsh: export GITHUB_TOKEN=\"$(gh auth token)\"\n" +
			"Then run: go test ./src/cmd/fetch ./src/e2e")
	}
	return token
}

func NewVersion(s string) models.Version {
	v, _ := models.ParseVersion(s)
	return v
}

func NewVersionPattern(s string) models.VersionPattern {
	v, _ := models.ParseVersionPattern(s)
	return v
}

func NewPRIdentifier(s string) models.PRIdentifier {
	prIdentifier, _ := models.ParsePrIdentifier(s)
	return prIdentifier
}
