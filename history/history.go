package history

import (
	"bufio"
	"os"
)

type History struct {
	Records Records
	Path    string
}

func Load(path string) (h *History, err error) {
	var records []Record
	h = &History{Records: Records{}}

	file, err := os.Open(path)
	if err != nil {
		// Return nil to regard it as no history (new)
		// if an open error occurs
		err = nil
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := &Record{}
		record.Unmarshal(scanner.Text())
		records = append(records, *record)
	}

	err = scanner.Err()
	if err != nil {
		return
	}

	return &History{
		Records: records,
		Path:    path,
	}, nil
}

func (h *History) Add(r Record) {
	h.Records = append(h.Records, r)
}

func (h *History) Save() error {
	file, err := os.OpenFile(h.Path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, record := range h.Records {
		b, _ := record.Marshal()
		w.Write(b)
		w.Write([]byte("\n"))
	}

	return w.Flush()
}
