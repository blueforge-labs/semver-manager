package tagprefix

import (
	"strings"

	"github.com/blueforge-labs/semver-manager/src/models"
)

// Strip removes prefix from a tag and reports whether the tag belongs to it.
func Strip(tag, prefix string) (string, bool) {
	if prefix == "" {
		return tag, true
	}
	if !strings.HasPrefix(tag, prefix) {
		return "", false
	}
	return strings.TrimPrefix(tag, prefix), true
}

// Parse converts matching tags or version lists into semantic versions.
func Parse(tags []string, prefix string) []models.Version {
	var versions []models.Version
	for _, tagList := range tags {
		for _, tag := range models.SplitVersions(tagList) {
			versionText, matches := Strip(tag, prefix)
			if !matches {
				continue
			}
			version, err := models.ParseVersion(versionText)
			if err == nil {
				versions = append(versions, version)
			}
		}
	}
	return versions
}

// Format returns a semantic version with the given tag prefix.
func Format(version models.Version, prefix string) string {
	return prefix + version.String()
}

// FormatVersions joins versions with spaces and adds prefix to each tag.
func FormatVersions(versions []models.Version, prefix string) string {
	formatted := make([]string, len(versions))
	for i, version := range versions {
		formatted[i] = Format(version, prefix)
	}
	return strings.Join(formatted, " ")
}
