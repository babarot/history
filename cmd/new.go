package cmd

import (
	"bufio"
	"os"

	"github.com/b4b4r07/history/cli"
	"github.com/najeira/ltsv"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [FILE/DIR]",
	Short: "",
	Long:  "",
	RunE:  new,
}

func new(cmd *cobra.Command, args []string) error {
	var err error
	data := []map[string]string{
		{"time": "05/Feb/2013:15:34:47 +0000", "host": "192.168.50.1", "req": "GET / HTTP/1.1", "status": "200"},
		{"time": "05/Feb/2013:15:35:15 +0000", "host": "192.168.50.1", "req": "GET /foo HTTP/1.1", "status": "200"},
		{"time": "05/Feb/2013:15:35:54 +0000", "host": "192.168.50.1", "req": "GET /bar HTTP/1.1", "status": "404"},
	}

	file, err := os.OpenFile(cli.Conf.History.Path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	w := ltsv.NewWriter(bufio.NewWriter(file))

	if err = w.WriteAll(data); err != nil {
		return err
	}
	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(newCmd)
}
