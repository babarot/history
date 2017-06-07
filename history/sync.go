package history

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/b4b4r07/history/config"
	"github.com/briandowns/spinner"
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

type Diff struct {
	Type    string
	Content string
}

func (h *History) Compare() (d Diff, err error) {
	fi, err := os.Stat(h.Path)
	if err != nil {
		err = errors.New("history file not found")
		return
	}

	ctx := context.Background()
	gist, _, err := h.client.Gists.Get(ctx, config.Conf.History.Sync.ID)
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
			err = fmt.Errorf("%s: not found on cloud", filepath.Base(h.Path))
			return
		}
		remoteContent = *file.Content
	}
	if remoteContent == localContent {
		d = Diff{Type: "equal", Content: ""}
		return
	}

	local := fi.ModTime().UTC()
	remote := gist.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return Diff{Type: "local", Content: localContent}, nil
	case remote.After(local):
		return Diff{Type: "remote", Content: remoteContent}, nil
	default:
	}
	d = Diff{Type: "equal", Content: ""}
	return
}

func (h *History) updateLocal(content string) error {
	return ioutil.WriteFile(h.Path, []byte(content), os.ModePerm)
}

func (h *History) updateRemote(content string) error {
	gist := github.Gist{
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(filepath.Base(h.Path)): {
				Content: github.String(content),
			},
		},
	}
	_, _, err := h.client.Gists.Edit(context.Background(), config.Conf.History.Sync.ID, &gist)
	return err
}

func (h *History) Sync() (err error) {
	var msg string
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Start()
	defer func() {
		if len(msg) > 0 {
			fmt.Println(msg)
		}
	}()
	defer s.Stop()

	h.client, err = getClient()
	if err != nil {
		return
	}

	if config.Conf.History.Sync.Token == "" {
		return errors.New("config history.sync.token is missing")
	}

	if config.Conf.History.Sync.ID == "" {
		return errors.New("config history.sync.id is missing")
	}

	diff, err := h.Compare()
	if err != nil {
		return
	}

	switch diff.Type {
	case "local":
		msg = "Uploaded"
		return h.updateRemote(diff.Content)
	case "remote":
		msg = "Downloaded"
		return h.updateLocal(diff.Content)
	case "equal":
		// Do nothing
	case "":
		// Locally but not remote
	default:
	}

	return
}
