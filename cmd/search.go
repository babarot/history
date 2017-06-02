package cmd

import (
	"fmt"

	"github.com/b4b4r07/history/cli"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [FILE/DIR]",
	Short: "",
	Long:  "",
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
	fmt.Printf("%#v\n", lines)
	return nil
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
