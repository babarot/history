package history

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/b4b4r07/history/config"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getClient() (gc *github.Client, err error) {
	cfg := config.Conf.History.Sync
	if cfg.Token == "" {
		err = errors.New("token is missing")
		return
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc), nil
}

func (h *History) Compare() (kind, content string, err error) {
	fi, err := os.Stat(h.Path)
	if err != nil {
		err = errors.New("history file not found")
		return
	}

	client, err := getClient()
	if err != nil {
		return
	}

	ctx := context.Background()
	gist, _, err := client.Gists.Get(ctx, config.Conf.History.Sync.ID)
	if err != nil {
		return
	}
	var (
		remoteContent, localContent string
	)
	out, err := ioutil.ReadFile(h.Path)
	if err != nil {
		return
	}
	localContent = string(out)
	for _, file := range gist.Files {
		if *file.Filename != filepath.Base(h.Path) {
			return "", "", fmt.Errorf("%s: not found on cloud", filepath.Base(h.Path))
		}
		remoteContent = *file.Content
	}
	if remoteContent == localContent {
		return "equal", "", nil
	}

	local := fi.ModTime().UTC()
	remote := gist.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return "local", localContent, nil
	case remote.After(local):
		return "remote", remoteContent, nil
	default:
	}

	return "equal", "", nil
}

func (h *History) updateLocal(content string) error {
	return ioutil.WriteFile(h.Path, []byte(content), os.ModePerm)
}

func (h *History) updateRemote(content string) error {
	gist := github.Gist{
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(filepath.Base(h.Path)): github.GistFile{
				Content: github.String(content),
			},
		},
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	_, _, err = client.Gists.Edit(context.Background(), config.Conf.History.Sync.ID, &gist)
	return err
}

func (h *History) Sync() (err error) {
	kind, content, err := h.Compare()
	if err != nil {
		return
	}
	log.Print(kind)
	switch kind {
	case "local":
		return h.updateRemote(content)
	case "remote":
		return h.updateLocal(content)
	case "equal":
		// Do nothing
	case "":
		// Locally but not remote
	default:
	}
	return
}
