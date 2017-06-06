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
	h, err := history.Load(config.Conf.History.Path)
	if err != nil {
		os.Exit(1)
	}

	h.Records.Sort()
	h.Records.Reverse()
	h.Records.Unique()
	h.Records.Reverse()

	if listBranch {
		h.Records.Branch(cli.GetBranchName())
	}

	if listDir {
		h.Records.Dir(cli.GetDirName())
	}

	if listQuery != "" {
		h.Records.Contains(listQuery)
	}

	for _, record := range h.Records {
		if listColumns == "" {
			fmt.Println(record.Raw())
		} else {
			// TODO
			config.Conf.History.Record.Columns = strings.Split(listColumns, ",")
			fmt.Println(record.Render())
		}
	}
}

var (
	listDir     bool
	listBranch  bool
	listQuery   string
	listColumns string
)

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listDir, "dir", "d", false, "List with dir")
	listCmd.Flags().BoolVarP(&listBranch, "branch", "b", false, "List with branch")
	listCmd.Flags().StringVarP(&listQuery, "query", "q", "", "List with query")
	listCmd.Flags().StringVarP(&listColumns, "columns", "c", "", "Specify columns with options")
}
