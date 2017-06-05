package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/b4b4r07/history/cli"
	"github.com/b4b4r07/history/config"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:   "config",
	Short: "Config the setting file",
	Long:  "Config the setting file with your editor (default: vim)",
	RunE:  conf,
}

var (
	confGetKey  string
	confAllKeys bool
)

func conf(cmd *cobra.Command, args []string) error {
	tomlfile := config.Conf.Core.TomlFile
	if tomlfile == "" {
		dir, _ := config.GetDefaultDir()
		tomlfile = filepath.Join(dir, "config.toml")
	}

	toml, err := toml.LoadFile(tomlfile)
	if err != nil {
		return err
	}

	if confAllKeys {
		allMap := toml.ToMap()
		for _, key := range toml.Keys() {
			fmt.Println(strings.Join(findKey(allMap, key), "\n"))
		}
		return nil
	}
	if confGetKey != "" {
		value := toml.Get(confGetKey)
		if value != nil {
			fmt.Printf("%v\n", value)
			return nil
		}
		return fmt.Errorf("%s: no such key found", confGetKey)
	}

	editor := config.Conf.Core.Editor
	if editor == "" {
		return cli.ErrConfigEditor
	}
	return cli.Run(editor, tomlfile)
}

func findKey(m map[string]interface{}, k string) []string {
	var ret []string
	originKey := k
	if v, ok := m[k]; ok {
		switch v.(type) {
		case map[string]interface{}:
			m = v.(map[string]interface{})
		default:
		}
	} else {
		return []string{}
	}
	for k, _ := range m {
		ret = append(ret, originKey+"."+k)
	}
	return ret
}

func init() {
	RootCmd.AddCommand(confCmd)
	confCmd.Flags().StringVarP(&confGetKey, "get", "", "", "Get the config value")
	confCmd.Flags().BoolVarP(&confAllKeys, "keys", "", false, "Get the config keys")
}
