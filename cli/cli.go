package cli

import (
	"errors"
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
