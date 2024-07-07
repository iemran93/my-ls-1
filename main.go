package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileInfo struct {
	os.FileInfo
	Path string
	Name string
}

type Path struct {
	Path    string
	Name    string
	Ptype   string
	LinkSrc string
}

type ByModTime []FileInfo

func (a ByModTime) Len() int {
	return len(a)
}

func (a ByModTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByModTime) Less(i, j int) bool {
	return a[i].ModTime().After(a[j].ModTime())
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = append(args, ".")
	}

	flags, inpPaths := parseArgs(args)
	// order paths to start with .

	for _, inpPath := range inpPaths {
		// get the path (Path struct)
		path, err := getPath(inpPath)
		if err != nil {
			log.Fatal(err)
		}

		if flags['R'] && path.Ptype == "Dir" {
			// use recursive printing if -R flag is set
			printRecursive(path, flags)
		} else if path.Ptype == "File" {
			// print the file if it is a file
			fileinfo, _ := os.Lstat(path.Path)
			files := []FileInfo{{FileInfo: fileinfo, Path: path.Path, Name: path.Name}}
			if len(inpPaths) > 1 {
				fmt.Printf("%s: \n", path.Name)
			}
			printFiles(path, files, flags)
			fmt.Print("\n")
		} else {
			// symlink
			absPath, _ := filepath.Abs(path.Path)
			fileinfo, _ := os.Lstat(absPath)
			if fileinfo.Mode()&os.ModeSymlink != 0 { // Check if it's a symlink
				if !flags['l'] {
					path.Path = absPath + "/"
				}
			}

			files, err := listDir(path, flags)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// handle multiple paths printing format
			if len(inpPaths) > 1 {
				fmt.Printf("%s: \n", path.Name)
			}
			printFiles(path, files, flags)
			fmt.Print("\n")
		}
	}
}
