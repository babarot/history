package cli

import (
	"errors"

	"github.com/b4b4r07/history/config"
)

var (
	ErrConfigEditor      = errors.New("config core.editor not set")
	ErrConfigHistoryPath = errors.New("config history.path not set")
)

func Edit(fname string) error {
	editor := config.Conf.Core.Editor
	if editor == "" {
		return ErrConfigEditor
	}

	return Run(editor, fname)
}
