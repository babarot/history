package cli

import (
	"errors"

	"github.com/b4b4r07/history/config"
)

var (
	ErrConfigEditor = errors.New("config editor not set")
)

func Edit(fname string) error {
	editor := config.Conf.Core.Editor
	if editor == "" {
		return ErrConfigEditor
	}

	return Run(editor, fname)
}
