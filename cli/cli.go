package cli

import (
	"errors"
	"strings"
)

var (
	ErrConfigEditor = errors.New("config editor not set")
)

func Edit(fname string) error {
	editor := Conf.Core.Editor
	if editor == "" {
		return ErrConfigEditor
	}

	return Run(editor, fname)
}

func IgnoringWord(command string) bool {
	command = strings.Split(command, " ")[0]
	for _, ignore := range Conf.History.Ignores {
		if ignore == command {
			return true
		}
	}
	return false
}

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}
