package push

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/go-github/v93/github"
)

// Config identifies the GitHub repository and optional local checkout.
type Config struct {
	Owner      string
	Repository string
	Token      string
	LocalPath  string
}

type gitAPI interface {
	GetRef(context.Context, string, string, string) (*github.Reference, *github.Response, error)
	CreateRef(context.Context, string, string, github.CreateRef) (*github.Reference, *github.Response, error)
}

var commitSHA = regexp.MustCompile(`^(?:[0-9a-fA-F]{40}|[0-9a-fA-F]{64})$`)

// PushTag creates a Git tag ref without creating a GitHub Release.
// target may be a commit SHA or GitHub ref; an empty target uses local HEAD.
func PushTag(ctx context.Context, config Config, tagName, target string) error {
	if config.Owner == "" || config.Repository == "" {
		return errors.New("owner and repository are required")
	}
	if config.Token == "" {
		return errors.New("GitHub token is required")
	}
	if tagName == "" {
		return errors.New("tag name is required")
	}

	client, err := github.NewClient(github.WithAuthToken(config.Token))
	if err != nil {
		return fmt.Errorf("create GitHub client: %w", err)
	}
	return pushTag(ctx, client.Git, config, tagName, target)
}

func pushTag(ctx context.Context, client gitAPI, config Config, tagName, target string) error {
	sha, err := resolveTarget(ctx, client, config, target)
	if err != nil {
		return err
	}
	_, _, err = client.CreateRef(ctx, config.Owner, config.Repository, github.CreateRef{
		Ref: "refs/tags/" + tagName,
		SHA: sha,
	})
	if err != nil {
		return fmt.Errorf("create tag %q: %w", tagName, err)
	}
	return nil
}

func resolveTarget(ctx context.Context, client gitAPI, config Config, target string) (string, error) {
	if target == "" {
		command := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
		if config.LocalPath != "" {
			command.Dir = config.LocalPath
		}
		output, err := command.Output()
		if err != nil {
			return "", fmt.Errorf("resolve local HEAD: %w", err)
		}
		return strings.TrimSpace(string(output)), nil
	}
	if commitSHA.MatchString(target) {
		return target, nil
	}

	ref := strings.TrimPrefix(target, "refs/")
	if !strings.HasPrefix(ref, "heads/") && !strings.HasPrefix(ref, "tags/") {
		ref = "heads/" + ref
	}
	reference, _, err := client.GetRef(ctx, config.Owner, config.Repository, ref)
	if err != nil {
		return "", fmt.Errorf("resolve GitHub ref %q: %w", target, err)
	}
	if reference == nil || reference.Object == nil || reference.Object.GetSHA() == "" {
		return "", fmt.Errorf("GitHub ref %q has no commit SHA", target)
	}
	return reference.Object.GetSHA(), nil
}
