package history

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func (h *History) Merge(a, b string) {
	lines := strings.Split(a+b, "\n")
	var records Records
	for _, line := range lines {
		if line == "" {
			continue
		}
		var r Record
		r.Unmarshal(line)
		records = append(records, r)
	}
	records.Sort()
	rs := make(Records, 0)
	encountered := map[Record]bool{}
	for _, record := range records {
		if !encountered[record] {
			encountered[record] = true
			rs = append(rs, record)
		}
	}
	h.Records = rs
}

func (h *History) updateLocal() error {
	var b bytes.Buffer
	for _, record := range h.Records {
		line, _ := record.Marshal()
		b.Write(line)
		b.WriteString("\n")
	}
	return ioutil.WriteFile(h.Path, b.Bytes(), os.ModePerm)
}

func (h *History) updateRemote() error {
	var b bytes.Buffer
	for _, record := range h.Records {
		line, _ := record.Marshal()
		b.Write(line)
		b.WriteString("\n")
	}
	gist := github.Gist{
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(filepath.Base(h.Path)): {
				Content: github.String(b.String()),
			},
		},
	}
	_, _, err := h.client.Gists.Edit(context.Background(), config.Conf.History.Sync.ID, &gist)
	return err
}

func (h *History) getGistID() (id string, err error) {
	var items []*github.Gist
	ctx := context.Background()

	// List items from gist.github.com
	gists, resp, err := h.client.Gists.List(ctx, "", &github.GistListOptions{})
	if err != nil {
		return
	}
	items = append(items, gists...)

	// pagenation
	for i := 2; i <= resp.LastPage; i++ {
		gists, _, err := h.client.Gists.List(ctx, "", &github.GistListOptions{
			ListOptions: github.ListOptions{Page: i},
		})
		if err != nil {
			continue
		}
		items = append(items, gists...)
	}

	for _, item := range items {
		for _, file := range item.Files {
			if *file.Filename == filepath.Base(config.Conf.History.Path) {
				id = *item.ID
				break
			}
		}
	}

	return
}

func (h *History) sync() (err error) {
	gist, _, err := h.client.Gists.Get(context.Background(), config.Conf.History.Sync.ID)
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

	h.Merge(remoteContent, localContent)
	if err := h.updateLocal(); err != nil {
		return err
	}
	if err := h.updateRemote(); err != nil {
		return err
	}

	return
}

func (h *History) Sync() (err error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Start()
	defer func() {
		if err == nil {
			fmt.Fprintln(os.Stderr, "Synced!")
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
		id, err := h.getGistID()
		if err != nil {
			return err
		}
		if id != "" {
			config.Conf.History.Sync.ID = id
		}
		if err := config.Conf.Save(); err != nil {
			return err
		}
	}

	if err := h.Backup(); err != nil {
		return err
	}

	return h.sync()
}
