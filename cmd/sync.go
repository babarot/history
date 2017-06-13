package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/b4b4r07/go-ask"
	"github.com/b4b4r07/history/config"
	"github.com/b4b4r07/history/history"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the history file with gist",
	Long:  "Sync the history file with gist",
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "\r"
	s.Writer = os.Stdout
	s.Start()
	defer s.Stop()

	h, err := history.Load()
	if err != nil {
		return err
	}
	if syncInterval > 0 {
		if skipSyncFor(syncInterval) {
			return fmt.Errorf("interval %v has not passed yet", syncInterval)
		}
	}
	diff, err := h.GetDiff()
	if err != nil {
		return err
	}
	if config.Conf.History.Sync.Size != 0 {
		if diff.Size < config.Conf.History.Sync.Size {
			return fmt.Errorf("The history difference %d is less than the specified size %d",
				diff.Size, config.Conf.History.Sync.Size)
		}
	}
	if syncAsk {
		s.Stop()
		if !ask.NewQ().Confirmf("%s: sync immediately?", config.Conf.History.Path) {
			return errors.New("canceled")
		}
	}
	s.Start()
	return h.Sync(diff)
}

func skipSyncFor(interval time.Duration) bool {
	file := filepath.Join(filepath.Dir(config.Conf.Core.TomlFile), ".sync")
	f, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		// Doesn't skip if some errors occur
		return false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		// Doesn't skip if some errors occur
		return false
	}
	if time.Now().Sub(fi.ModTime()).Hours() < interval.Hours() {
		// Skip if the fixed time has not elapsed
		return true
	}
	// Update the timestamp if sync
	os.Chtimes(file, time.Now(), time.Now())
	return false
}

var (
	syncInterval time.Duration
	syncAsk      bool
)

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().DurationVarP(&syncInterval, "interval", "", 0, "Sync with the interval")
	syncCmd.Flags().BoolVarP(&syncAsk, "ask", "", false, "Sync after the confirmation")
	syncCmd.Flags().IntVarP(&config.Conf.History.Sync.Size, "diff", "", config.Conf.History.Sync.Size, "Sync if the diff exceeds a certain number")
}
