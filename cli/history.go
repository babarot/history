package cli

import (
	"bufio"
	"os"
	"time"

	ltsv "github.com/Songmu/go-ltsv"
	"github.com/google/uuid"
)

type History struct {
	ID   uint32
	Date time.Time

	Command string
	Dir     string
	Branch  string
	Status  int
}

func NewHistory() *History {
	return &History{
		ID:   uuid.New().ID(),
		Date: time.Now(),
	}
}

func (h *History) SetCommand(arg string) { h.Command = arg }
func (h *History) SetDir(arg string)     { h.Dir = arg }
func (h *History) SetBranch(arg string)  { h.Branch = arg }
func (h *History) SetStatus(arg int)     { h.Status = arg }

func (h *History) Add() error {
	file, err := os.OpenFile(Conf.History.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ltsv.Marshal(h)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	w.Write(b)
	w.Write([]byte("\n"))

	return w.Flush()
}
