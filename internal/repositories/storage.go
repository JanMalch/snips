package repositories

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"snips/internal/utils"
	"strconv"
	"strings"

	yaml "github.com/goccy/go-yaml"
)

var (
	ErrInvalidRepositoryDefinition = errors.New("invalid repository definition")
	ErrCannotFindRepository        = errors.New("cannot find repository")
	ErrDuplicateRepositoryName     = errors.New("multiple repositories with the same name")
)

type storage struct {
	Repositories []storedRepository `yaml:"repositories"`
}

type storedGitHubRepository struct {
	Repo   string `yaml:"repo"`
	Owner  string `yaml:"owner"`
	Branch string `yaml:"branch,omitempty"`
	Prefix string `yaml:"prefix,omitempty"`
}

type storedRepository struct {
	Name      string                  `yaml:"name,omitempty"`
	LocalPath string                  `yaml:"path,omitempty"`
	Url       string                  `yaml:"url,omitempty"`
	GitHub    *storedGitHubRepository `yaml:"github,omitempty"`
}

func readStorage() (storage, error) {
	f, err := utils.RepositoriesFile()
	if err != nil {
		return storage{}, err
	}
	b, err := os.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			return storage{
				Repositories: make([]storedRepository, 0),
			}, nil
		}
		return storage{}, err
	}
	var stored storage
	err = yaml.Unmarshal(b, &stored)
	if err != nil {
		return storage{}, err
	}
	return stored, nil
}

func validate(storage storage) error {
	seens := make(map[string]bool)
	for _, r := range storage.Repositories {
		_, seen := seens[r.Name]
		if seen {
			return ErrDuplicateRepositoryName
		}
		seens[r.Name] = true
	}
	return nil
}

func writeStorage(storage storage) error {
	err := validate(storage)
	if err != nil {
		return err
	}
	f, err := utils.RepositoriesFile()
	if err != nil {
		return err
	}
	b, err := yaml.Marshal(storage)
	if err != nil {
		return err
	}
	return os.WriteFile(f, b, 0644)
}

func ReadStored() ([]Repository, error) {
	stored, err := readStorage()
	if err != nil {
		return nil, err
	}
	repos := make([]Repository, 0)
	for _, storedRepo := range stored.Repositories {
		var repo Repository
		if storedRepo.LocalPath != "" {
			repo = Local(storedRepo.Name, storedRepo.LocalPath)
		} else if storedRepo.Url != "" {
			repo, err = Http(storedRepo.Name, storedRepo.Url)
			if err != nil {
				return nil, err
			}
		} else if storedRepo.GitHub != nil {
			// FIXME: auth token
			repo, err = GitHub(storedRepo.Name, storedRepo.GitHub.Owner, storedRepo.GitHub.Repo, storedRepo.GitHub.Branch, storedRepo.GitHub.Prefix, "")
			if err != nil {
				return nil, err
			}
		}
		if repo == nil {
			panic("unable to determine type for stored repository")
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func FindStored(q string) (Repository, error) {
	stored, err := ReadStored()
	if err != nil {
		return nil, err
	}
	if len(stored) == 0 {
		return nil, ErrCannotFindRepository
	}
	if q == "" {
		if len(stored) > 1 {
			return nil, ErrCannotFindRepository
		}
		return stored[0], nil
	}
	if index, err := strconv.Atoi(q); err == nil {
		if index >= len(stored) {
			return nil, ErrCannotFindRepository
		}
		if index < 0 {
			return stored[len(stored)+index], nil
		}
		return stored[index], nil
	}
	for _, s := range stored {
		if s.Name() == q {
			return s, nil
		}
	}
	return nil, ErrCannotFindRepository
}

func OverwriteStored(repos []Repository) error {
	storedRepos := make([]storedRepository, 0)
	for _, repo := range repos {
		switch v := repo.(type) {
		case *localDirRepository:
			storedRepos = append(storedRepos, storedRepository{
				Name:      v.name,
				LocalPath: v.dir,
			})
		case *httpRepository:
			storedRepos = append(storedRepos, storedRepository{
				Name: v.name,
				Url:  v.url,
			})
		case *gitHubRepository:
			storedRepos = append(storedRepos, storedRepository{
				Name: v.name,
				GitHub: &storedGitHubRepository{
					Repo:   v.repo,
					Owner:  v.owner,
					Branch: v.branch,
					Prefix: v.prefix,
				},
			})
		case nil:
			break
		default:
			panic(fmt.Sprintf("unknown repository type %T", repo))
		}
	}
	return writeStorage(storage{
		Repositories: storedRepos,
	})
}

func AddStored(name, input string) error {
	if name == "" || input == "" {
		return ErrInvalidRepositoryDefinition
	}
	storedRepo := storedRepository{
		Name: name,
	}
	if strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://") {
		parsedURL, err := url.Parse(input)
		if err != nil {
			return err
		}
		if parsedURL.Host == "github.com" {
			parts := strings.SplitN(parsedURL.Path[1:], "/", 5)
			switch len(parts) {
			case 0, 1:
				// fallback to https repo
				storedRepo.Url = input
			case 2, 3:
				storedRepo.GitHub = &storedGitHubRepository{
					Owner: parts[0],
					Repo:  parts[1],
				}
			case 4:
				storedRepo.GitHub = &storedGitHubRepository{
					Owner:  parts[0],
					Repo:   parts[1],
					Branch: parts[3],
				}
			default:
				storedRepo.GitHub = &storedGitHubRepository{
					Owner:  parts[0],
					Repo:   parts[1],
					Branch: parts[3],
					Prefix: parts[4],
				}
			}
		} else {
			storedRepo.Url = input
		}
	} else {
		abs, err := filepath.Abs(input)
		if err != nil {
			return err
		}
		storedRepo.LocalPath = abs
	}
	storage, err := readStorage()
	if err != nil {
		return err
	}
	storage.Repositories = append(storage.Repositories, storedRepo)
	return writeStorage(storage)
}
