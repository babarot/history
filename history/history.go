package history

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/b4b4r07/history/config"
	"github.com/google/go-github/github"
)

type History struct {
	Records Records
	Path    string

	client *github.Client
}

func Load(path string) (h *History, err error) {
	var records []Record
	h = &History{Records: Records{}, Path: path}

	file, err := os.Open(path)
	if err != nil {
		// Return nil to regard it as no history (new)
		// if an open error occurs
		err = nil
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := &Record{}
		record.Unmarshal(scanner.Text())
		records = append(records, *record)
	}

	err = scanner.Err()
	if err != nil {
		return
	}

	return &History{
		Records: records,
		Path:    path,
	}, nil
}

func (h *History) Save() error {
	file, err := os.OpenFile(h.Path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, record := range h.Records {
		b, _ := record.Marshal()
		w.Write(b)
		w.Write([]byte("\n"))
	}

	return w.Flush()
}

func (h *History) Backup() (err error) {
	if _, err := os.Stat(h.Path); err != nil {
		// cannot backup if no history
		return nil
	}

	dir, err := config.GetDefaultDir()
	if err != nil {
		return
	}

	dir = filepath.Join(dir, ".backup", time.Now().Format("2006/01/02"))
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return
	}

	src, err := os.Open(h.Path)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(dir, filepath.Base(h.Path)))
	if err != nil {
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
