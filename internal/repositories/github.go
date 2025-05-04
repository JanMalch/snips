package repositories

import (
	"context"
	"errors"
	"fmt"
	"snips/internal/snippets"
	"strings"

	"github.com/google/go-github/v71/github"
)

var (
	ErrGitHubNoBranch = errors.New("no branch given")
)

type gitHubRepository struct {
	name   string
	owner  string
	repo   string
	branch string
	prefix string
	client *github.Client
}

func GitHub(name, owner, repo, branch, prefix, authToken string) (Repository, error) {
	client := github.NewClient(nil)
	if authToken != "" {
		client = client.WithAuthToken(authToken)
	}
	if branch == "" {
		repo, _, err := client.Repositories.Get(context.Background(), owner, repo)
		if err != nil {
			return nil, err
		}
		if repo.DefaultBranch != nil && *repo.DefaultBranch != "" {
			branch = *repo.DefaultBranch
		} else {
			return nil, ErrGitHubNoBranch
		}
	}
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	return &gitHubRepository{
		name:   name,
		owner:  owner,
		repo:   repo,
		branch: branch,
		prefix: prefix,
		client: client,
	}, nil
}

func (r *gitHubRepository) joinDescriptiveParts() string {
	return fmt.Sprintf("github.com/%s/%s/tree/%s/%s", r.owner, r.repo, r.branch, r.prefix)
}

func (r *gitHubRepository) String() string {
	return r.name + " @ " + r.joinDescriptiveParts()
}

func (r *gitHubRepository) Name() string {
	return r.name
}

func (r *gitHubRepository) recurseListAll(path string) (snippets.Ids, error) {
	_, contents, _, err := r.client.Repositories.GetContents(context.Background(), r.owner, r.repo, path, &github.RepositoryContentGetOptions{
		Ref: r.branch,
	})
	if err != nil {
		return []snippets.Id{}, err
	}
	if contents == nil {
		return []snippets.Id{}, nil
	}
	ids := make([]snippets.Id, 0)
	for _, content := range contents {
		if content.Type == nil {
			continue
		}
		switch *content.Type {
		case "directory":
			children, err := r.recurseListAll(path + "/" + *content.Name)
			if err != nil {
				return nil, err
			}
			ids = append(ids, children...)
		case "file":
			id, err := snippets.NewId(strings.TrimPrefix(path+"/"+*content.Name, r.prefix))
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		default:
			panic("Unknown GitHub content type: " + *content.Type)
		}
	}
	return ids, nil
}

func (r *gitHubRepository) ListAll() (snippets.Ids, error) {
	ids, err := r.recurseListAll(r.prefix)
	if err != nil {
		return nil, err
	}
	snippets.Sort(ids)
	return ids, nil
}

func (r *gitHubRepository) Read(id snippets.Id) (string, error) {
	path := r.prefix + id.String()
	content, _, _, err := r.client.Repositories.GetContents(context.Background(), r.owner, r.repo, path, &github.RepositoryContentGetOptions{
		Ref: r.branch,
	})
	if content == nil || err != nil {
		return "", ErrNoSuchSnippet
	}
	return content.GetContent()
}
