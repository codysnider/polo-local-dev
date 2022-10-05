package git

import (
	"context"
	gogit "github.com/go-git/go-git/v5"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
	"os"
)

var (
	Client *github.Client
)

func init() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	Client = github.NewClient(tc)
}

func GetOrganization() (*github.Organization, error) {
	ctx := context.Background()
	org, _, err := Client.Organizations.Get(ctx, "poloniex")
	if err != nil {
		return nil, err
	}

	return org, nil
}

func GetOrganizationRepo(name string) (*github.Repository, error) {
	ctx := context.Background()
	repo, _, err := Client.Repositories.Get(ctx, "poloniex", name)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func ForkRepo(repo *github.Repository) error {
	fork, _, forkErr := Client.Repositories.CreateFork(context.Background(), "poloniex", repo.GetName(), &github.RepositoryCreateForkOptions{})
	if forkErr != nil && fork == nil {
		return forkErr
	}

	return nil
}

func CloneRepo(destination string, repo *github.Repository) error {
	_, cloneErr := gogit.PlainClone(destination, false, &gogit.CloneOptions{
		URL: repo.GetCloneURL(),
	})

	return cloneErr
}
