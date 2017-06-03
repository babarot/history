package cmd

import (
	"fmt"
	"os"

	"github.com/b4b4r07/history/cli"
	"github.com/b4b4r07/history/history"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long:  "list",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	h, err := history.Load(cli.Conf.History.Path)
	if err != nil {
		os.Exit(1)
	}

	h.Records.Sort()
	h.Records.Reverse()
	h.Records.Unique()
	h.Records.Reverse()
	h.Records.Grep(args)

	for _, record := range h.Records {
		fmt.Println(record.Command)
	}
}

func init() {
	RootCmd.AddCommand(listCmd)
}
