package push

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-github/v93/github"
	"github.com/stretchr/testify/require"
)

type fakeGitAPI struct {
	getRefName string
	getRef     *github.Reference
	createdRef github.CreateRef
}

func (f *fakeGitAPI) GetRef(_ context.Context, _, _, ref string) (*github.Reference, *github.Response, error) {
	f.getRefName = ref
	return f.getRef, nil, nil
}

func (f *fakeGitAPI) CreateRef(_ context.Context, _, _ string, ref github.CreateRef) (*github.Reference, *github.Response, error) {
	f.createdRef = ref
	return nil, nil, nil
}

func TestPushTagWithSHA(t *testing.T) {
	sha := strings.Repeat("a", 40)
	api := &fakeGitAPI{}
	err := pushTag(context.Background(), api, Config{Owner: "org", Repository: "repo"}, "release/1.2.3", sha)

	require.NoError(t, err)
	require.Equal(t, "refs/tags/release/1.2.3", api.createdRef.Ref)
	require.Equal(t, sha, api.createdRef.SHA)
}

func TestPushTagResolvesRemoteBranch(t *testing.T) {
	sha := strings.Repeat("b", 40)
	api := &fakeGitAPI{getRef: &github.Reference{Object: &github.GitObject{SHA: &sha}}}
	err := pushTag(context.Background(), api, Config{Owner: "org", Repository: "repo"}, "v1.2.3", "main")

	require.NoError(t, err)
	require.Equal(t, "heads/main", api.getRefName)
	require.Equal(t, sha, api.createdRef.SHA)
}

func TestPushTagDefaultsToLocalHead(t *testing.T) {
	repo := t.TempDir()
	for _, args := range [][]string{{"init", "-q", repo}, {"-C", repo, "-c", "user.name=test", "-c", "user.email=test@example.com", "commit", "--allow-empty", "-qm", "initial"}} {
		command := exec.Command("git", args...)
		output, err := command.CombinedOutput()
		require.NoError(t, err, string(output))
	}
	output, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
	require.NoError(t, err)
	sha := strings.TrimSpace(string(output))
	api := &fakeGitAPI{}

	err = pushTag(context.Background(), api, Config{Owner: "org", Repository: "repo", LocalPath: repo}, "1.2.3", "")
	require.NoError(t, err)
	require.Equal(t, sha, api.createdRef.SHA)
}
