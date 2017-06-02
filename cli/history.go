package cli

import (
	"bufio"
	"fmt"
	"os"
	"time"

	ltsv "github.com/Songmu/go-ltsv"
	"github.com/google/uuid"
)

type Record struct {
	ID      uint32
	Date    time.Time
	Command string
	Dir     string
	Branch  string
	Status  int
}

type Records []Record

type History struct {
	Records Records
}

func NewHistory() (h *History, err error) {
	var rs Records
	file, err := os.Open(Conf.History.Path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := Record{}
		ltsv.Unmarshal([]byte(scanner.Text()), &r)
		rs = append(rs, r)
	}

	err = scanner.Err()
	if err != nil {
		return
	}

	return &History{
		Records: rs,
	}, nil
}

func (h *History) Add(r Record) {
	h.Records = append(h.Records, r)
}

func (h *History) Save() error {
	file, err := os.OpenFile(Conf.History.Path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, record := range h.Records {
		b, err := ltsv.Marshal(record)
		if err != nil {
			return err
		}
		w.Write(b)
		w.Write([]byte("\n"))
	}

	return w.Flush()
}

func (r *Record) render() string {
	return fmt.Sprintf("%d\t%s", r.ID, r.Command)
}

func NewRecord() *Record {
	return &Record{
		ID:   uuid.New().ID(),
		Date: time.Now(),
	}
}

func (r *Record) SetCommand(arg string) { r.Command = arg }
func (r *Record) SetDir(arg string)     { r.Dir = arg }
func (r *Record) SetBranch(arg string)  { r.Branch = arg }
func (r *Record) SetStatus(arg int)     { r.Status = arg }
