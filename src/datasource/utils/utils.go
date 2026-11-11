package utils

import (
	"context"
	"fmt"

	"github.com/google/go-github/v93/github"
	"golang.org/x/oauth2"
)

type DatasourceClient interface {
	listTags(owner string, repo string) ([]string, error)
}

type DryDatasourceClient struct {
	Username   string
	Password   string
	Repository string
}

func (drySv *DryDatasourceClient) listTags(owner string, repo string) ([]string, error) {
	return []string{}, nil
}

type Datasource struct {
	client DatasourceClient
}

type Filters struct {
	Highest bool
	Release bool
}

type GithubClient struct {
	client *github.Client
	err    error
}

func newGithubClient(token string) *GithubClient {
	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client, err := github.NewClient(github.WithHTTPClient(oauthClient))

	return &GithubClient{
		client: client,
		err:    err,
	}
}

func (ghClient *GithubClient) listTags(owner string, repo string) ([]string, error) {
	if ghClient.err != nil {
		return nil, fmt.Errorf("create GitHub client: %w", ghClient.err)
	}
	tagList := []string{}
	ctx := context.Background()
	opts := &github.ListOptions{PerPage: 100, Page: 1}

	for {
		tags, resp, err := ghClient.client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return []string{}, fmt.Errorf("ListTags error: %w", err)
		}

		for _, tag := range tags {
			tagList = append(tagList, tag.GetName())
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return tagList, nil
}

func NewSemverSvc(platform, token string) (svc Datasource) {
	if platform == "github" {
		svc = Datasource{client: newGithubClient(token)}
	} else if platform == "dry-run" {
		svc = Datasource{client: &DryDatasourceClient{}}
	}

	return svc
}

func (s *Datasource) FetchSemverTags(owner, repo string) (tagList []string, err error) {
	if s.client == nil {
		return nil, fmt.Errorf("git platform client is not defined")
	}

	tagList, err = s.client.listTags(owner, repo)
	if err != nil {
		return nil, err
	}

	return tagList, err
}
