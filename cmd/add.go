package cmd

import (
	"github.com/b4b4r07/history/cli"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new history",
	Long:  "Add new history",
	RunE:  add,
}

func add(cmd *cobra.Command, args []string) error {
	h := cli.NewHistory()
	h.SetCommand("git status")
	h.SetDir("/Users/b4b4r07/src/github.com/b4b4r07/gist")
	return h.Add()
}

func init() {
	RootCmd.AddCommand(addCmd)
}
