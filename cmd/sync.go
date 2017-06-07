package cmd

import (
	"github.com/b4b4r07/history/config"
	"github.com/b4b4r07/history/history"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the history file with gist",
	Long:  "Sync the history file with gist",
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) error {
	h, err := history.Load(config.Conf.History.Path)
	if err != nil {
		return err
	}

	// Add record to history
	return h.Sync()
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
