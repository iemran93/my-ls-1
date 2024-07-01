package main

import (
	"log"
	"strings"
)

func parseArgs(args []string) (map[rune]bool, []string) {
	flags := map[rune]bool{}
	var paths []string
	flagsList := []rune{'a', 'l', 'r', 'R', 't'}

	for _, arg := range args {
		if arg == "-" {
			paths = append(paths, arg)
			continue
		}
		if len(arg) >= 2 && arg[:2] == "--" {
			if arg == "--" {
				continue
			}
			if arg == "--reverse" {
				flags[rune('r')] = true
			} else if arg == "--recursive" {
				flags[rune('R')] = true
			} else if arg == "--all" {
				flags[rune('a')] = true
			} else {
				log.Fatalf("option '%s' is ambiguous; possibilities: '--all' '--reverse' '--recursive'", arg)
			}
		} else if strings.HasPrefix(arg, "-") {
			for _, flag := range arg[1:] {
				if !strings.ContainsRune(string(flagsList), flag) {
					log.Fatalf("ls: invalid option -- '%c'", flag)
				}
				flags[rune(flag)] = true
			}
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	return flags, paths
}
