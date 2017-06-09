package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/b4b4r07/history/cli"
	"github.com/b4b4r07/history/config"
	"github.com/b4b4r07/history/history"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the history",
	Long:  "List the history",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	h, err := history.Load()
	if err != nil {
		os.Exit(1)
	}

	h.Records.Sort()
	h.Records.Reverse()
	h.Records.Unique()
	h.Records.Reverse()

	if config.Conf.Screen.FilterDir {
		h.Records.Dir(cli.GetDirName())
	}
	if config.Conf.Screen.FilterBranch {
		h.Records.Branch(cli.GetBranchName())
	}
	if config.Conf.Screen.FilterHostname {
		config.Conf.Var.Hostname = cli.GetHostName()
	}
	if config.Conf.Var.Query != "" {
		h.Records.Contains(config.Conf.Var.Query)
	}

	for _, record := range h.Records {
		if config.Conf.Var.Columns == "" {
			fmt.Println(record.Raw())
		} else {
			// TODO
			config.Conf.Screen.Columns = strings.Split(config.Conf.Var.Columns, ",")
			fmt.Println(record.Render())
		}
	}
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&config.Conf.Screen.FilterDir, "filter-dir", "d", config.Conf.Screen.FilterDir, "List with dir")
	listCmd.Flags().BoolVarP(&config.Conf.Screen.FilterBranch, "filter-branch", "b", config.Conf.Screen.FilterBranch, "List with branch")
	listCmd.Flags().BoolVarP(&config.Conf.Screen.FilterHostname, "filter-hostname", "p", config.Conf.Screen.FilterHostname, "List with hostname")
	listCmd.Flags().StringVarP(&config.Conf.Var.Query, "query", "q", config.Conf.Var.Query, "List with query")
	listCmd.Flags().StringVarP(&config.Conf.Var.Columns, "columns", "c", config.Conf.Var.Columns, "Specify columns with options")
}
