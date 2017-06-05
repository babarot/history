package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	Path     string       `toml:"path"`
	Ignores  []string     `toml:"ignores"`
	Sync     SyncConfig   `toml:"sync"`
	Record   RecordConfig `toml:"record"`
	UseColor bool         `toml:"use_color"`
}

type SyncConfig struct{}

type RecordConfig struct {
	Visible  []string `toml:"visible"`
	StatusOK string   `toml:"status_ok"`
	StatusNG string   `toml:"status_ng"`
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
	cfg.History.UseColor = false
	cfg.History.Record.Visible = []string{"{{.Command}}"}
	cfg.History.Record.StatusOK = "o"
	cfg.History.Record.StatusNG = "x"

	return toml.NewEncoder(f).Encode(cfg)
}

func CheckIgnores(command string) bool {
	command = strings.Split(command, " ")[0]
	for _, ignore := range Conf.History.Ignores {
		if ignore == command {
			return true
		}
	}
	return false
}

func KeyCol(vs []string) int {
	for i, v := range vs {
		if v == "{{.Command}}" {
			return i
		}
	}
	return -1
}
