package cli

import (
	"bytes"
	"errors"
	"strings"

	"github.com/b4b4r07/history/history"
)

type Screen struct {
	Lines   []string
	Records []history.Record
}

func NewScreen(c ScreenConfig) (s *Screen, err error) {
	var (
		lines   []string
		records history.Records
	)

	h, err := history.Load(Conf.History.Path)
	if err != nil {
		return
	}

	h.Records.Sort()
	h.Records.Reverse()
	h.Records.Unique()

	for _, record := range h.Records {
		if c.Dir != "" && c.Dir != record.Dir {
			continue
		}
		if c.Branch != "" && c.Branch != record.Branch {
			continue
		}
		lines = append(lines, record.Render(Conf.History.Visible))
		records = append(records, record)
	}

	return &Screen{
		Lines:   lines,
		Records: records,
	}, nil
}

type Line struct {
	history.Record
}

type Lines []Line

func (s *Screen) parseLine(line string) (*Line, error) {
	l := strings.Split(line, "\t")
	var record history.Record
	idx := Index(Conf.History.Visible, "{{.Command}}")
	if idx > len(l) {
		return &Line{}, errors.New("invalid index; review config visible")
	}
	for _, record = range s.Records {
		if record.Command == l[idx] {
			break
		}
	}
	return &Line{record}, nil
}

func (l *Lines) Filter(fn func(Line) bool) *Lines {
	lines := make(Lines, 0)
	for _, line := range *l {
		if fn(line) {
			lines = append(lines, line)
		}
	}
	return &lines
}

func (s *Screen) Select() (lines Lines, err error) {
	if len(s.Lines) == 0 {
		err = errors.New("no text to display")
		return
	}
	selectcmd := Conf.Core.SelectCmd
	if selectcmd == "" {
		err = errors.New("no selectcmd specified")
		return
	}

	text := strings.NewReader(strings.Join(s.Lines, "\n"))
	var buf bytes.Buffer
	err = runFilter(selectcmd, text, &buf)
	if err != nil {
		return
	}

	if buf.Len() == 0 {
		err = errors.New("no lines selected")
		return
	}

	selectedLines := strings.Split(buf.String(), "\n")
	for _, line := range selectedLines {
		if line == "" {
			continue
		}
		parsedLine, err := s.parseLine(line)
		if err != nil {
			// TODO: log
			continue
		}
		lines = append(lines, *parsedLine)
	}

	if len(lines) == 0 {
		err = errors.New("no lines selected")
		return
	}

	return
}
