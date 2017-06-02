package cli

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"

	ltsv "github.com/Songmu/go-ltsv"
)

type Screen struct {
	Lines []string
}

func NewScreen() (s *Screen, err error) {
	var lines []string

	file, err := os.Open(Conf.History.Path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		h := History{}
		ltsv.Unmarshal([]byte(scanner.Text()), &h)
		// lines = append(lines, h.Command)
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return
	}

	return &Screen{
		Lines: lines,
	}, nil
}

type Line struct {
	ID          string
	ShortID     string
	Description string
	Filename    string
	Path        string
	URL         string
	Public      bool
}

type Lines []Line

func (s *Screen) parseLine(line string) (*Line, error) {
	return &Line{
		ID: line,
	}, nil
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
