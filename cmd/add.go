package cmd

import (
	"errors"

	"github.com/b4b4r07/history/config"
	"github.com/b4b4r07/history/history"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new history",
	Long:  "Add new history",
	RunE:  add,
}

func add(cmd *cobra.Command, args []string) error {
	h, err := history.Load(config.Conf.History.Path)
	if err != nil {
		return err
	}

	r := history.NewRecord()
	if addCommand == "" {
		return errors.New("--command option is required")
	}
	if addDir == "" {
		return errors.New("--dir option is required")
	}

	// Skip adding if the command is registed as ignoring word
	if config.CheckIgnores(addCommand) {
		return nil
	}

	r.SetCommand(addCommand)
	r.SetDir(addDir)
	r.SetBranch(addBranch)
	r.SetStatus(addStatus)

	// Backup before adding new record
	if err := h.Backup(); err != nil {
		return err
	}
	h.Records.Add(*r)

	return h.Save()
}

var (
	addCommand string
	addDir     string
	addBranch  string
	addStatus  int
)

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addCommand, "command", "", "", "Set command")
	addCmd.Flags().StringVarP(&addDir, "dir", "", "", "Set dir")
	addCmd.Flags().StringVarP(&addBranch, "branch", "", "", "Set branch")
	addCmd.Flags().IntVarP(&addStatus, "status", "", 0, "Set status")
}
