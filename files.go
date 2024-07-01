package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func getPath(inpPath string) (Path, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Path{}, err
	}

	wdParent := filepath.Dir(wd)
	inpPath = strings.Replace(inpPath, wdParent, "..", -1)
	inpPath = strings.Replace(inpPath, wd, ".", -1)

	path := ""
	name := ""
	if len(inpPath) >= 2 && inpPath[:2] == ".." {
		after := strings.TrimLeft(inpPath, "..")
		path = filepath.Dir(wd) + after
		name = ".." + after
	} else if inpPath[0] == '.' {
		after := strings.TrimLeft(inpPath, ".")
		path = wd + after
		name = "." + after
	} else {
		path = inpPath
		name = inpPath
	}

	fileinfo, err := os.Stat(path)
	if err != nil {
		return Path{}, err
	}

	pType := ""
	if fileinfo.IsDir() {
		pType = "Dir"
	} else {
		pType = "File"
	}

	nPath := Path{
		Path:  path,
		Name:  name,
		Ptype: pType,
	}

	return nPath, nil
}

func listDir(dir Path, flags map[rune]bool) ([]FileInfo, error) {
	var files []FileInfo
	skipfirst := true

	// add "." and ".." explicitly if -a flag is set
	if flags['a'] {
		currentDir := dir.Path
		parentDir := filepath.Dir(dir.Path)

		currentStat, err := os.Stat(currentDir)
		if err != nil {
			return nil, err
		}
		files = append(files, FileInfo{FileInfo: currentStat, Path: dir.Path, Name: "."})

		parentStat, err := os.Stat(parentDir)
		if err != nil {
			return nil, err
		}
		files = append(files, FileInfo{FileInfo: parentStat, Path: parentDir, Name: ".."})
	}

	err := filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip the root directory
		if skipfirst {
			skipfirst = false
			return nil
		}

		// a flag not set
		if !flags['a'] && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Base(filepath.Dir(path)) != filepath.Base(dir.Path) {
			return filepath.SkipDir
		}

		files = append(files, FileInfo{FileInfo: info, Path: path, Name: info.Name()})

		return nil
	})

	if err != nil {
		return nil, err
	}

	if flags['a'] {
		sortFilesCustom(files)
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
