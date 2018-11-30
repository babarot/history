package config

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

type Path struct {
	path string
}

func NewPath(path string) Path {
	p := Path{path: path}
	return p
}

func (p *Path) UnmarshalText(text []byte) error {
	p.path = string(text)
	return nil
}

func (p *Path) MarshalText() (text []byte, err error) {
	return []byte(p.path), nil
}

func (p *Path) Abs() string {
	path := p.path

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "$HOME/") {
		home := ""
		if home = os.Getenv("HOME"); home == "" {
			user, err := user.Current()
			if err != nil {
				log.Fatalf("Failed to get user home and $HOME not set.")
			}
			home = user.HomeDir
		}

		if home == "" {
			log.Fatalf("Failed to get user home and $HOME not set.")
		}

		if strings.HasPrefix(path, "~/") {
			path = strings.Replace(path, "~/", home+"/", 1)
		}
		if strings.HasPrefix(path, "$HOME/") {
			path = strings.Replace(path, "$HOME/", home+"/", 1)
		}
	}
	return path
}

type Config struct {
	Core    CoreConfig    `toml:"core"`
	History HistoryConfig `toml:"history"`
	Screen  ScreenConfig  `toml:"screen"`

	// Var cooperates with other packages
	Var VarConfig `toml:"-"`
}

type CoreConfig struct {
	Editor    string `toml:"editor"`
	SelectCmd string `toml:"selectcmd"`
	TomlFile  Path   `toml:"tomlfile"`
}

type HistoryConfig struct {
	Path       Path       `toml:"path"`
	BackupPath Path       `toml:"backup_path"`
	Ignores    []string   `toml:"ignore_words"`
	Sync       SyncConfig `toml:"sync"`
	UseColor   bool       `toml:"use_color"`
}

type SyncConfig struct {
	ID    string `toml:"id"`
	Token string `toml:"token"`
	Size  int    `toml:"size"`
}

type ScreenConfig struct {
	FilterDir      bool     `toml:"filter_dir"`
	FilterBranch   bool     `toml:"filter_branch"`
	FilterHostname bool     `toml:"filter_hostname"`
	FilterStatus   bool     `toml:"filter_status"`
	Columns        []string `toml:"columns"`
	StatusOK       string   `toml:"status_ok"`
	StatusNG       string   `toml:"status_ng"`
}

type VarConfig struct {
	Dir      string
	Branch   string
	Hostname string
	Status 	 bool
	Query    string
	Columns  string
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

func (cfg *Config) Save() error {
	f, err := os.OpenFile(cfg.Core.TomlFile.Abs(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return toml.NewEncoder(f).Encode(cfg)
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
	cfg.Core.TomlFile = NewPath(file)

	cfg.History.Path = NewPath(filepath.Join(dir, "history.ltsv"))
	cfg.History.BackupPath = NewPath(filepath.Join(dir, ".backup"))

	cfg.History.Ignores = []string{}
	cfg.History.UseColor = false
	cfg.History.Sync.ID = ""
	cfg.History.Sync.Token = "$GITHUB_TOKEN"
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		cfg.History.Sync.Token = token
	}
	cfg.History.Sync.Size = 100

	cfg.Screen.FilterDir = false
	cfg.Screen.FilterBranch = false
	cfg.Screen.FilterHostname = false
	cfg.Screen.FilterStatus = false
	cfg.Screen.Columns = []string{"{{.Time}}", "{{.Status}}", "{{.Command}}"}
	cfg.Screen.StatusOK = " "
	cfg.Screen.StatusNG = "x"

	return toml.NewEncoder(f).Encode(cfg)
}
