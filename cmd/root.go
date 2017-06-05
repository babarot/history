package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/b4b4r07/history/config"
	"github.com/spf13/cobra"
)

const Version = "0.0.1"

var showVersion bool

var RootCmd = &cobra.Command{
	Use:           "history",
	Short:         "Enhanced shell history with LTSV",
	Long:          "Enhanced shell history with LTSV",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("version %s/%s\n", Version, runtime.Version())
			return
		}
		cmd.Usage()
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	initConf()
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show the version and exit")
}

func initConf() {
	dir, _ := config.GetDefaultDir()
	toml := filepath.Join(dir, "config.toml")

	err := config.Conf.LoadFile(toml)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
