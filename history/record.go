package history

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	tt "text/template"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	ltsv "github.com/Songmu/go-ltsv"
	pathshorten "github.com/b4b4r07/go-pathshorten"
	"github.com/b4b4r07/history/config"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	pipeline "github.com/mattn/go-pipeline"
)

type Record struct {
	Date     time.Time
	Command  string
	Dir      string
	Branch   string
	Status   int
	Hostname string
}

type Records []Record

func NewRecord() *Record {
	hostname, _ := os.Hostname()
	return &Record{
		Date:     time.Now(),
		Hostname: hostname,
	}
}

func (r *Record) SetCommand(arg string) { r.Command = arg }
func (r *Record) SetDir(arg string)     { r.Dir = arg }
func (r *Record) SetBranch(arg string)  { r.Branch = arg }
func (r *Record) SetStatus(arg int)     { r.Status = arg }

func (r *Record) Raw() string {
	out, _ := r.Marshal()
	return string(out)
}

func (r *Record) Render() (line string) {
	var tmpl *tt.Template
	columns := config.Conf.History.Record.Columns
	if len(columns) == 0 {
		// default
		columns = []string{"{{.Command}}"}
	}
	format := columns[0]
	for _, v := range columns[1:] {
		format += "\t" + v
	}
	t, err := tt.New("format").Parse(format)
	if err != nil {
		return
	}
	tmpl = t
	if tmpl != nil {
		var b bytes.Buffer
		err := tmpl.Execute(&b, map[string]interface{}{
			"Date":    r.Date.Format("2006-01-02"),
			"Time":    fmt.Sprintf("%-15s", humanize.Time(r.Date)),
			"Command": r.renderCommand(),
			"Dir":     r.renderDir(),
			"Path":    r.Dir,
			"Base":    color.BlueString(filepath.Base(r.Dir)),
			"Branch":  r.Branch,
			"Status": func(status int) string {
				switch status {
				case 0:
					ok := config.Conf.History.Record.StatusOK
					if ok == "" {
						ok = "o"
					}
					return color.GreenString(ok)
				default:
					ng := config.Conf.History.Record.StatusNG
					if ng == "" {
						ng = "x"
					}
					return color.RedString(ng)
				}
			}(r.Status),
			"Hostname": r.Hostname,
		})
		if err != nil {
			return
		}
		line = b.String()
	}
	return
}

func (r *Record) renderCommand() string {
	if !config.Conf.History.UseColor {
		return r.Command
	}
	highlight, err := exec.LookPath("highlight")
	if err != nil {
		return r.Command
	}
	// TODO: more faster
	out, err := pipeline.Output(
		[]string{"echo", r.Command},
		[]string{highlight, "-S", "sh", "-O", "ansi"},
	)
	if err != nil {
		return r.Command
	}
	return strings.TrimSuffix(string(out), "\n")
}

func (r *Record) renderDir() string {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 20
	}
	dir := r.Dir
	if len(dir) > w/3 {
		dir = pathshorten.Run(dir)
	}
	return color.BlueString(dir)
}

func (r *Record) Unmarshal(line string) {
	ltsv.Unmarshal([]byte(line), r)
}

func (r *Record) Marshal() (line []byte, err error) {
	b, err := ltsv.Marshal(r)
	if err != nil {
		return
	}
	return b, nil
}

func (rs *Records) Add(r Record) {
	*rs = append(*rs, r)
}

func (rs *Records) Delete(r Record) {
	*rs = *rs.Reduce(func(rr Record) bool {
		return rr.Command == r.Command && rr.Dir == r.Dir && rr.Branch == r.Branch
	})
}

func (r *Records) Reduce(fn func(Record) bool) *Records {
	records := make(Records, 0)
	for _, record := range *r {
		if !fn(record) {
			records = append(records, record)
		}
	}
	return &records
}

func (r *Records) Filter(fn func(Record) bool) *Records {
	records := make(Records, 0)
	for _, record := range *r {
		if fn(record) {
			records = append(records, record)
		}
	}
	return &records
}

func (r *Records) Unique() {
	rs := make(Records, 0)
	encountered := map[string]bool{}
	for _, record := range *r {
		if !encountered[record.Command] {
			encountered[record.Command] = true
			rs = append(rs, record)
		}
	}
	*r = rs
}

func (r *Records) Reverse() {
	var rs Records
	for i := len(*r) - 1; i >= 0; i-- {
		rs = append(rs, (*r)[i])
	}
	*r = rs
}

func (r *Records) Contains(word string) {
	*r = *r.Filter(func(r Record) bool {
		return strings.Contains(r.Command, word)
	})
}

func (r *Records) Branch(branch string) {
	*r = *r.Filter(func(r Record) bool {
		return r.Branch == branch
	})
}

func (r *Records) Dir(dir string) {
	*r = *r.Filter(func(r Record) bool {
		return r.Dir == dir
	})
}

func (r *Records) Latest() Record {
	if len(*r) < 1 {
		return Record{}
	}
	return (*r)[len(*r)-1]
}

func (r Records) Len() int           { return len(r) }
func (r Records) Less(i, j int) bool { return r[i].Date.Before(r[j].Date) }
func (r Records) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r *Records) Sort() {
	sort.Sort(*r)
}

func init() {
	color.NoColor = false
}
