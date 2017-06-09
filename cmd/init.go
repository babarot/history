package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/b4b4r07/history/config"
	"github.com/b4b4r07/history/history"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Migrate your ~/.zsh_history at first",
	Long:  "Migrate your ~/.zsh_history at first",
	RunE:  migration,
}

var re = regexp.MustCompile(`^: (\d+):\d+;(.*)$`)

func migration(cmd *cobra.Command, args []string) error {
	if !strings.Contains(os.Getenv("SHELL"), "zsh") {
		return errors.New("Supports zsh only now")
	}

	h, err := history.Load()
	if err != nil {
		return err
	}

	if len(h.Records) > 10 {
		return fmt.Errorf("%s: there are already many histories", config.Conf.History.Path)
	}

	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".zsh_history"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		columns := re.FindStringSubmatch(scanner.Text())
		if len(columns) != 3 {
			continue
		}
		timestamp := columns[1]
		command := columns[2]
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return err
		}
		r := history.NewRecord()
		r.SetDate(time.Unix(ts, 0))
		r.SetCommand(command)
		h.Records.Add(*r)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return h.Save()
}

func init() {
	RootCmd.AddCommand(initCmd)
}
