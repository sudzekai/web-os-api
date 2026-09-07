package args

import (
	"slices"
	"strings"
)

var Args = make(map[string]string)

var Flags []string = []string{
	"--dev",
	"--cfg",
}

func LoadArguments(args []string) {
	for _, flag := range Flags {
		idx := slices.Index(args, flag)
		if idx != -1 {
			Args[flag] = ""

			if len(args) >= idx+2 {
				value := args[idx+1]
				if !strings.HasPrefix(value, "--") {
					Args[flag] = args[idx+1]
				}
			}
		}
	}
}
