package history

import (
	"testing"

	"github.com/b4b4r07/history/config"
)

func TestCheckIgnores(t *testing.T) {
	config.Conf.History.Ignores = []string{
		`^cd(\s+-+\w?)?$`,
	}
	tests := []struct {
		command string
		want    bool
	}{
		{command: "cd", want: true},
		{command: "cd ", want: false},
		{command: " cd", want: false},
		{command: "cd -", want: true},
		{command: "cd -G", want: true},
	}
	for _, test := range tests {
		got := CheckIgnores(test.command)
		if got != test.want {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}
}

func TestIndexCommandColumns(t *testing.T) {
	tests := []struct {
		columns []string
		want    int
	}{
		{columns: []string{"{{.Command}}"}, want: 0},
		{columns: []string{"{{.Time}}", "{{.Status}}"}, want: -1},
		{columns: []string{"{{.Time}}", "{{.Status}}", "{{.Command}}"}, want: 2},
	}
	for _, test := range tests {
		config.Conf.Screen.Columns = test.columns
		got := IndexCommandColumns()
		if got != test.want {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}
}
