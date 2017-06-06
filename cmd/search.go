package cmd

import (
	"fmt"

	"github.com/b4b4r07/history/cli"
	"github.com/b4b4r07/history/config"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the command from the history file",
	Long:  "Search the command from the history file",
	RunE:  search,
}

func search(cmd *cobra.Command, args []string) error {
	if config.Conf.Screen.FilterDir {
		config.Conf.Screen.Dir = cli.GetDirName()
	}
	if config.Conf.Screen.FilterBranch {
		config.Conf.Screen.Branch = cli.GetBranchName()
	}

	screen, err := cli.NewScreen()
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

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVarP(&config.Conf.Screen.FilterDir, "dir", "d", config.Conf.Screen.FilterDir, "Search with dir")
	searchCmd.Flags().BoolVarP(&config.Conf.Screen.FilterBranch, "branch", "b", config.Conf.Screen.FilterBranch, "Search with branch")
	searchCmd.Flags().StringVarP(&config.Conf.Screen.Query, "query", "q", config.Conf.Screen.Query, "Search with query")
	searchCmd.Flags().StringVarP(&config.Conf.Screen.Columns, "columns", "c", config.Conf.Screen.Columns, "Specify columns with options")
}
