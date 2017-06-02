package cmd

import (
	"fmt"

	"github.com/b4b4r07/history/cli"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the command from the history file",
	Long:  "Search the command from the history file",
	RunE:  search,
}

func search(cmd *cobra.Command, args []string) error {
	screen, err := cli.NewScreen()
	if err != nil {
		return err
	}

	lines, err := screen.Select()
	if err != nil {
		return err
	}

	command := lines[0].Command
	for _, line := range lines[1:] {
		command += "; " + line.Command
	}
	fmt.Println(command)

	return nil
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
