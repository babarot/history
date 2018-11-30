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
		config.Conf.Var.Dir = cli.GetDirName()
	}
	if config.Conf.Screen.FilterBranch {
		config.Conf.Var.Branch = cli.GetBranchName()
	}
	if config.Conf.Screen.FilterHostname {
		config.Conf.Var.Hostname = cli.GetHostName()
	}
	if config.Conf.Screen.FilterStatus {
		config.Conf.Var.Status = true
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
	searchCmd.Flags().BoolVarP(&config.Conf.Screen.FilterDir, "filter-dir", "d", config.Conf.Screen.FilterDir, "Search with dir")
	searchCmd.Flags().BoolVarP(&config.Conf.Screen.FilterBranch, "filter-branch", "b", config.Conf.Screen.FilterBranch, "Search with branch")
	searchCmd.Flags().BoolVarP(&config.Conf.Screen.FilterHostname, "filter-hostname", "p", config.Conf.Screen.FilterHostname, "Search with hostname")
	searchCmd.Flags().BoolVarP(&config.Conf.Screen.FilterStatus, "filter-status", "s", config.Conf.Screen.FilterStatus, "Search with status OK")
	searchCmd.Flags().StringVarP(&config.Conf.Var.Query, "query", "q", config.Conf.Var.Query, "Search with query")
	searchCmd.Flags().StringVarP(&config.Conf.Var.Columns, "columns", "c", config.Conf.Var.Columns, "Specify columns with options")
}
