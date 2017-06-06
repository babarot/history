package cmd

import (
	"github.com/b4b4r07/history/cli"
	"github.com/b4b4r07/history/config"
	"github.com/b4b4r07/history/history"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete the command from the history file",
	Long:  "delete the command from the history file",
	RunE:  delete,
}

func delete(cmd *cobra.Command, args []string) error {
	screen, err := cli.NewScreen(deleteConfig())
	if err != nil {
		return err
	}

	lines, err := screen.Select()
	if err != nil {
		return err
	}

	h, err := history.Load(config.Conf.History.Path)
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

func deleteConfig() config.ScreenConfig {
	cfg := config.ScreenConfig{}
	if deleteDir {
		cfg.Dir = cli.GetDirName()
	}
	if deleteBranch {
		cfg.Branch = cli.GetBranchName()
	}
	if deleteQuery != "" {
		cfg.Query = deleteQuery
	}
	if deleteColumns != "" {
		cfg.Columns = deleteColumns
	}
	return cfg
}

var (
	deleteDir     bool
	deleteBranch  bool
	deleteQuery   string
	deleteColumns string
)

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deleteDir, "dir", "d", false, "delete with dir")
	deleteCmd.Flags().BoolVarP(&deleteBranch, "branch", "b", false, "delete with branch")
	deleteCmd.Flags().StringVarP(&deleteQuery, "query", "q", "", "delete with query")
	deleteCmd.Flags().StringVarP(&deleteColumns, "columns", "c", "", "Specify columns with options")
}
