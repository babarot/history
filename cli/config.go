package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Core    CoreConfig    `toml:"core"`
	History HistoryConfig `toml:"history"`
}

type CoreConfig struct {
	Editor    string `toml:"editor"`
	SelectCmd string `toml:"selectcmd"`
	TomlFile  string `toml:"tomlfile"`
}

type HistoryConfig struct {
	Path    string       `toml:"path"`
	Ignores []string     `toml:"ignores"`
	Sync    SyncConfig   `toml:"sync"`
	Record  RecordConfig `toml:"record"`
}

type SyncConfig struct{}

type RecordConfig struct {
	Visible []string `toml:"visible"`
}

// ScreenConfig is only for Screen
type ScreenConfig struct {
	Dir    string
	Branch string
	Query  string
}

var Conf Config

func GetDefaultDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	}
	dir = filepath.Join(dir, "history")

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return dir, fmt.Errorf("cannot create directory: %v", err)
	}

	return dir, nil
}

func (cfg *Config) LoadFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	// base dir
	dir := filepath.Dir(file)

	cfg.Core.Editor = os.Getenv("EDITOR")
	if cfg.Core.Editor == "" {
		cfg.Core.Editor = "vim"
	}
	cfg.Core.SelectCmd = "fzf-tmux --multi:fzf --multi:peco"
	cfg.Core.TomlFile = file

	cfg.History.Path = filepath.Join(dir, "history.ltsv")
	cfg.History.Ignores = []string{}
	cfg.History.Record.Visible = []string{"{{.Command}}"}

	return toml.NewEncoder(f).Encode(cfg)
}
