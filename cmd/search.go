package cmd

import (
	"fmt"

	"github.com/b4b4r07/history/cli"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the command from the history file",
	Long:  "Search the command from the history file",
	RunE:  search,
}

func search(cmd *cobra.Command, args []string) error {
	screen, err := cli.NewScreen(searchConfig())
	if err != nil {
		return err
	}

	lines, err := screen.Select()
	if err != nil {
		return err
	}

	command := lines[0].Command
	for _, line := range lines[1:] {
		command += "; " + line.Command
	}
	fmt.Println(command)

	return nil
}

func searchConfig() cli.ScreenConfig {
	cfg := cli.ScreenConfig{}
	if searchDir {
		cfg.Dir = cli.GetDirName()
	}
	if searchBranch {
		cfg.Branch = cli.GetBranchName()
	}
	return cfg
}

var (
	searchDir    bool
	searchBranch bool
)

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVarP(&searchDir, "dir", "d", false, "Search with dir")
	searchCmd.Flags().BoolVarP(&searchBranch, "branch", "b", false, "Search with branch")
}
