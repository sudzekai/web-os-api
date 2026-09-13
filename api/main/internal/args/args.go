package args

import (
	"slices"
	"strings"
)

var Args = make(map[string]string)

var flags []string = []string{
	"--dev",
	"--cfg",
	"--cli",
	"--log-level",
}

func LoadArguments(args []string) {
	for _, flag := range flags {
		idx := slices.Index(args, flag)
		if idx != -1 {
			Args[flag] = "val"

			if len(args) >= idx+2 {
				value := args[idx+1]
				if !strings.HasPrefix(value, "--") {
					Args[flag] = args[idx+1]
				}
			}
		}
	}
}
