package config

import "strings"

func CheckIgnores(command string) bool {
	command = strings.Split(command, " ")[0]
	for _, ignore := range Conf.History.Ignores {
		if ignore == command {
			return true
		}
	}
	return false
}

func IndexCommandColumns() int {
	for i, v := range Conf.Screen.Columns {
		if v == "{{.Command}}" {
			return i
		}
	}
	return -1
}
