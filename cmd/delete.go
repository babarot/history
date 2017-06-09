package cmd

import (
	"github.com/b4b4r07/history/cli"
	"github.com/b4b4r07/history/config"
	"github.com/b4b4r07/history/history"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the record from history file",
	Long:  "Delete the selected record from history file",
	RunE:  delete,
}

func delete(cmd *cobra.Command, args []string) error {
	if config.Conf.Screen.FilterDir {
		config.Conf.Var.Dir = cli.GetDirName()
	}
	if config.Conf.Screen.FilterBranch {
		config.Conf.Var.Branch = cli.GetBranchName()
	}
	if config.Conf.Screen.FilterHostname {
		config.Conf.Var.Hostname = cli.GetHostName()
	}

	screen, err := cli.NewScreen()
	if err != nil {
		return err
	}

	lines, err := screen.Select()
	if err != nil {
		return err
	}

	h, err := history.Load()
	if err != nil {
		return err
	}

	for _, line := range lines {
		h.Records.Delete(history.Record{
			Command: line.Command,
			Dir:     line.Dir,
			Branch:  line.Branch,
		})
	}

	return h.Save()
}

var (
	deleteDir, deleteBranch bool
)

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&config.Conf.Screen.FilterDir, "filter-dir", "d", config.Conf.Screen.FilterDir, "Delete with dir")
	deleteCmd.Flags().BoolVarP(&config.Conf.Screen.FilterBranch, "filter-branch", "b", config.Conf.Screen.FilterBranch, "Delete with branch")
	deleteCmd.Flags().BoolVarP(&config.Conf.Screen.FilterHostname, "filter-hostname", "p", config.Conf.Screen.FilterHostname, "Delete with hostname")
	deleteCmd.Flags().StringVarP(&config.Conf.Var.Query, "query", "q", config.Conf.Var.Query, "Delete with query")
	deleteCmd.Flags().StringVarP(&config.Conf.Var.Columns, "columns", "c", config.Conf.Var.Columns, "Specify columns with options")
}
