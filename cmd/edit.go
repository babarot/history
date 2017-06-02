package cmd

import (
	"github.com/b4b4r07/history/cli"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit your history file directly",
	Long:  "Edit your history file directly",
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	return cli.Edit(cli.Conf.History.Path)
}

func init() {
	RootCmd.AddCommand(editCmd)
}
