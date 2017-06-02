package cli

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/b4b4r07/history/history"
)

type Screen struct {
	Lines   []string
	Records []history.Record
}

func NewScreen() (s *Screen, err error) {
	var (
		lines   []string
		records history.Records
	)

	file, err := os.Open(Conf.History.Path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := &history.Record{}
		r.Unmarshal(scanner.Text())
		lines = append(lines, r.Render())
		records = append(records, *r)
	}

	err = scanner.Err()
	if err != nil {
		return
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
	id, _ := strconv.Atoi(l[0])
	var r history.Record
	for _, record := range s.Records {
		if record.ID == uint32(id) {
			r = record
		}
	}
	return &Line{
		r,
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
