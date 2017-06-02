package history

import (
	"fmt"
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

func (r *Record) Render() string {
	return fmt.Sprintf("%d\t%s\t%s", r.ID, r.Date.Format("2006-01-02"), r.Command)
}

func (r *Record) Unmarshal(line string) Record {
	ltsv.Unmarshal([]byte(line), r)
	return *r
}

func (r *Record) Marshal() ([]byte, error) {
	b, err := ltsv.Marshal(r)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
