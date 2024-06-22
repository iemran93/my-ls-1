//go:build !unix
// +build !unix

package main

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

type FileInfo struct {
	os.FileInfo
	Path string
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

	flags, paths := parseArgs(args)

	for _, path := range paths {
		files, err := listDir(path, flags)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		printFiles(path, files, flags)
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
	skipfirst := true
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if skipfirst {
			skipfirst = false
			return nil
		}

		// a flag not set
		/* if !flags['a'] && strings.HasPrefix(path, ".") {
			return nil
		} */
		// log.Println(info.Name(), filepath.Base(dir))
		// log.Println(filepath.Dir(path), filepath.Base(dir), info.Name())
		if !flags['a'] && strings.HasPrefix(info.Name(), ".") && path != "." {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// R flag not set
		/* 		if info.IsDir() && path != "." && !flags['R'] {
			return filepath.SkipDir
		} */
		/* if !flags['R'] && strings.Contains(path, "/") {
			return filepath.SkipDir
		} */
		if !flags['R'] && filepath.Base(filepath.Dir(path)) != filepath.Base(dir) {
			return filepath.SkipDir
		}

		files = append(files, FileInfo{info, path})

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

func printFiles(dir string, files []FileInfo, flags map[rune]bool) {
	if flags['R'] {
		// using dircsKeys to maintaine the dircs map order
		dircs, dircsKeys := prepareRecursive(dir, files) // map[string][]FileInfo
		for _, dircKey := range dircsKeys {
			k := dircKey
			v := dircs[dircKey]
			fmt.Printf("%s:\n", k)
			// print size
			if flags['l'] {
				printDirSize(v)
			}
			for indx, file := range v {
				if file.Path == "" { // skip nil(empty directory)
					fmt.Print("")
					continue
				}
				if flags['l'] {
					printLongFormat(file)
				} else {
					if indx != len(v)-1 {
						if file.IsDir() {
							cFolder := colorFolder(file)
							fmt.Print(cFolder, " ")
						} else {
							fmt.Print(file.Name(), " ")
						}
					} else {
						if file.IsDir() {
							cFolder := colorFolder(file)
							fmt.Printf("%s\n", cFolder)
						} else {
							fmt.Printf("%s\n", file.Name())
						}
					}
				}
			}
			fmt.Println()
		}
	} else {
		if flags['l'] {
			printDirSize(files)
		}
		for indx, file := range files {
			if flags['l'] {
				printLongFormat(file)
			} else {
				if indx != len(files)-1 {
					if file.IsDir() {
						cFolder := colorFolder(file)
						fmt.Print(cFolder, " ")
					} else {
						fmt.Print(file.Name(), " ")
					}
				} else {
					if file.IsDir() {
						cFolder := colorFolder(file)
						fmt.Printf("%s\n", cFolder)
					} else {
						fmt.Printf("%s\n", file.Name())
					}
				}
			}
		}
	}
}

func prepareRecursive(dir string, files []FileInfo) (map[string][]FileInfo, []string) {
	dircs := map[string][]FileInfo{}
	var dircsKeys []string
	// dirBase := filepath.Dir(dir)
	for _, file := range files {
		fileBase := filepath.Dir(file.Path)
		if fileBase != dir {
			fileBase = fmt.Sprintf(filepath.Dir(dir) + "/" + fileBase)
		}
		if _, exist := dircs[fileBase]; !exist {
			dircsKeys = append(dircsKeys, fileBase)
		}
		dircs[fileBase] = append(dircs[fileBase], file)
		if file.FileInfo.IsDir() {
			contents, _ := os.ReadDir(file.Path)
			if len(contents) == 0 {
				path := fmt.Sprint(filepath.Dir(dir) + "/" + file.Path)
				dircs[path] = append(dircs[path], FileInfo{})
				dircsKeys = append(dircsKeys, path)
			}
		}
	}
	// 	log.Println(file)
	// 	fileBase := filepath.Dir(file.Path)

	// 	if fileBase != dirBase {
	// 		// if _, exists := dircs[fileBase]; !exists {
	// 		// 	dircsKeys = append(dircsKeys, fileBase)
	// 		// }

	// 		if file.IsDir() {
	// 			if _, exist := dircs[fileBase]; !exist {
	// 				dircs[fileBase] = append(dircs[fileBase], FileInfo{})
	// 				dircsKeys = append(dircsKeys, fileBase)
	// 			} else {
	// 				dircs[fileBase] = append(dircs[fileBase], file)
	// 			}
	// 		} else {
	// 			dircs[fileBase] = append(dircs[fileBase], file)
	// 		}
	// 	} else {
	// 		// if _, exists := dircs[dirBase]; !exists {
	// 		// 	dircsKeys = append(dircsKeys, dirBase)
	// 		// }

	// 		if file.FileInfo.IsDir() {
	// 			if _, exist := dircs[dirBase]; !exist {
	// 				dircs[dirBase] = append(dircs[dirBase], FileInfo{})
	// 				dircsKeys = append(dircsKeys, dirBase)
	// 			} else {
	// 				dircs[dirBase] = append(dircs[dirBase], file)
	// 			}

	// 		} else {
	// 			dircs[dirBase] = append(dircs[dirBase], file)
	// 		}
	// 	}
	// }
	// log.Println(dircs, dircsKeys)
	// // handle empty folder

	// // remove nil from non-empty folders
	// for k, v := range dircs {
	// 	if len(v) > 1 {
	// 		newSlice := []FileInfo{}
	// 		for _, file := range v {
	// 			if file.Path != "" {
	// 				newSlice = append(newSlice, file)
	// 			}
	// 		}
	// 		dircs[k] = newSlice
	// 	}
	// }
	return dircs, dircsKeys
}

func printLongFormat(file FileInfo) {

	modTime := file.ModTime().Format(time.Stamp)
	modTime = modTime[:len(modTime)-3]

	cFolder := colorFolder(file)

	uName, gName := getUserGroup(file)
	nOfLinks := getLinksNumber(file)

	fmt.Printf("%s\t%d\t%s\t%s\t%d\t%s\t",
		file.Mode().String(),
		nOfLinks, // Placeholder for Nlink
		uName,    // Placeholder for Username
		gName,    // Placeholder for Group
		file.Size(),
		modTime)
	if file.IsDir() {
		fmt.Printf("%s\n", cFolder)
	} else {
		fmt.Printf("%s\n",
			file.Name())
	}
}

func printDirSize(files []FileInfo) {
	var dirSize int
	if files[0].Path != "" {
		for _, file := range files {
			if file.Size() < 1024 {
				dirSize += 1
			} else {
				dirSize += int(math.Ceil(float64(file.Size()) / 1024))
			}
		}
	}
	// convert byte to MB
	fmt.Printf("total %d\n", dirSize)
}

func getUserGroup(file FileInfo) (string, string) {
	uName, gName := "N/A", "N/A"

	var uid, gid uint32
	if stat, ok := file.Sys().(*syscall.Stat_t); ok {
		uid = stat.Uid
		gid = stat.Gid
	} else {
		return uName, gName
	}

	userLN, err := user.LookupId(fmt.Sprint(uid))
	if err == nil {
		uName = userLN.Username
	}

	grouplG, err := user.LookupGroupId(fmt.Sprint(gid))
	if err == nil {
		gName = grouplG.Name
	}

	return uName, gName
}

func getLinksNumber(file FileInfo) uint64 {
	var nOfLinks uint64
	if stat, ok := file.Sys().(*syscall.Stat_t); ok {
		return stat.Nlink
	} else {
		return nOfLinks
	}
}

func colorFolder(file FileInfo) string {
	colorReset := "\033[0m"
	colorBlueText := "\033[34m"
	colorLightGreenBackground := "\033[102m"

	coloredFolder := fmt.Sprintf("%s%s%s%s", colorBlueText,
		colorLightGreenBackground, file.Name(), colorReset)

	return coloredFolder
}
