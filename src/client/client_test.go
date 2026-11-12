package client

import (
	"testing"

	"github.com/blueforge-labs/semver-manager/src/models"
	"github.com/blueforge-labs/semver-manager/src/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestPrefixFormatsAndParsesReleaseTags(t *testing.T) {
	c := New(Config{Prefix: "release/"})
	tags := c.parseTags([]string{"release/1.2.3", "release/2.0.0-rc.1", "v1.0.0", "release/not-semver"})

	require.Len(t, tags, 2)
	require.Equal(t, "1.2.3", tags[0].String())
	require.Equal(t, "release/2.0.0-rc.1", c.FormatRelease(tags[1]))
}

func TestValidatePrefixedVersion(t *testing.T) {
	c := New(Config{Prefix: "v"})

	valid, err := c.Validate("v1.2.3", "strict")
	require.NoError(t, err)
	require.True(t, valid)

	valid, err = c.Validate("1.2.3", "strict")
	require.Error(t, err)
	require.False(t, valid)
}

func TestFilterAndIncrement(t *testing.T) {
	c := New(Config{})
	versions, err := models.ParseVersions("1.0.0 1.0.1")
	require.NoError(t, err)

	filtered, err := c.Filter(versions, filter.Highest())
	require.NoError(t, err)
	require.Equal(t, "1.0.1", filtered[0].String())

	stream, err := models.ParseVersionPattern("1.0.*")
	require.NoError(t, err)
	incremented, err := c.Increment(versions, stream, models.Patch)
	require.NoError(t, err)
	require.Equal(t, "1.0.2", incremented.String())
}
