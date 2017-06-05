package cmd

import (
	"fmt"
	"os"

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

	for _, arg := range args {
		h.Records.Contains(arg)
	}

	for _, record := range h.Records {
		fmt.Println(record.Command)
	}
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listDir, "dir", "d", false, "List with dir")
	listCmd.Flags().BoolVarP(&listBranch, "branch", "b", false, "List with branch")
	listCmd.Flags().StringVarP(&listQuery, "query", "q", "", "List with query")
}

var (
	listDir    bool
	listBranch bool
	listQuery  string
)
