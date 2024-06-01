// +build !unix

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileInfo struct {
	os.FileInfo
	Path string
}

type ByModTime []FileInfo

func (a ByModTime) Len() int           { return len(a) }
func (a ByModTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByModTime) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = append(args, ".")
	}

	flags, paths := parseArgs(args)

	for _, path := range paths {
		files, err := listDir(path, flags)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		printFiles(files, flags)
	}
}

func parseArgs(args []string) (map[rune]bool, []string) {
	flags := make(map[rune]bool)
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, flag := range arg[1:] {
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

func listDir(dir string, flags map[rune]bool) ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !flags['a'] && strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		files = append(files, FileInfo{info, path})

		if !flags['R'] && path != dir {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if flags['t'] {
		sort.Sort(ByModTime(files))
	}

	if flags['r'] {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	return files, nil
}

func printFiles(files []FileInfo, flags map[rune]bool) {
	for _, file := range files {
		if flags['l'] {
			printLongFormat(file)
		} else {
			fmt.Println(file.Path)
		}
	}
}

func printLongFormat(file FileInfo) {
	modTime := file.ModTime().Format(time.RFC3339)
	fmt.Printf("%s\t%d\t%s\t%s\t%d\t%s\t%s\n",
		file.Mode().String(),
		0, // Placeholder for Nlink
		"unknown", // Placeholder for Username
		"unknown", // Placeholder for Group
		file.Size(),
		modTime,
		file.Path)
}