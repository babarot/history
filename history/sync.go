package history

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/b4b4r07/history/config"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getClient() (gc *github.Client, err error) {
	cfg := config.Conf.History.Sync
	if cfg.Token == "" {
		err = errors.New("config history.sync.token is missing")
		return
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: expandPath(cfg.Token)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc), nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
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
			if *file.Filename == filepath.Base(config.Conf.History.Path.Abs()) {
				id = *item.ID
				break
			}
		}
	}

	// Case that couldn't be found
	if id == "" {
		id, err = h.create()
	}

	return
}

func (h *History) create() (id string, err error) {
	out, err := ioutil.ReadFile(h.Path)
	if err != nil {
		return
	}
	localContent := string(out)
	var (
		files = map[github.GistFilename]github.GistFile{
			github.GistFilename(filepath.Base(h.Path)): {
				Content: github.String(localContent),
			},
		}
		public = false
		desc   = ""
	)
	item, resp, err := h.client.Gists.Create(context.Background(), &github.Gist{
		Files:       files,
		Public:      &public,
		Description: &desc,
	})
	if item == nil {
		err = errors.New("unexpected fatal error: item is nil")
		return
	}
	if resp == nil {
		err = errors.New("Try again when you have a better network connection")
		return
	}
	id = *item.ID
	return
}

type Diff struct {
	Local struct {
		Size    int
		Content string
	}
	Remote struct {
		Size    int
		Content string
	}
	Size int
}

func (h *History) GetDiff() (d Diff, err error) {
	h.client, err = getClient()
	if err != nil {
		return
	}

	if config.Conf.History.Sync.ID == "" {
		id, err := h.getGistID()
		if err != nil {
			return d, err
		}
		if id != "" {
			config.Conf.History.Sync.ID = id
		}
		if err := config.Conf.Save(); err != nil {
			return d, err
		}
	}

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

	return Diff{
		Local: struct {
			Size    int
			Content string
		}{
			Size:    strings.Count(localContent, "\n"),
			Content: localContent,
		},
		Remote: struct {
			Size    int
			Content string
		}{
			Size:    strings.Count(remoteContent, "\n"),
			Content: remoteContent,
		},
		Size: int(math.Abs(float64(
			strings.Count(localContent, "\n") - strings.Count(remoteContent, "\n"),
		))),
	}, nil
}

func (h *History) Sync(diff Diff) (err error) {
	if err := h.Backup(); err != nil {
		return err
	}

	h.Merge(diff.Remote.Content, diff.Local.Content)
	if err := h.updateLocal(); err != nil {
		return err
	}
	if err := h.updateRemote(); err != nil {
		return err
	}

	return
}
