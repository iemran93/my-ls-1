package main

import (
	"fmt"
	"log"
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
	Name string
}

type Path struct {
	Path  string
	Name  string
	Ptype string
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
	log.Println(flags, inpPaths)
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
			fileinfo, _ := os.Stat(path.Path)
			files := []FileInfo{{FileInfo: fileinfo, Path: path.Path, Name: path.Name}}
			if len(inpPaths) > 1 {
				fmt.Printf("%s: \n", path.Name)
			}
			printFiles(path, files, flags)
			fmt.Print("\n")
		} else {
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

func printFiles(dir Path, files []FileInfo, flags map[rune]bool) {
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
					fmt.Print(file.Name, " ")
				}
			} else {
				if file.IsDir() {
					cFolder := colorFolder(file)
					fmt.Printf("%s\n", cFolder)
				} else {
					fmt.Printf("%s\n", file.Name)
				}
			}
		}
	}
}

func printRecursive(dir Path, flags map[rune]bool) {
	files, err := listDir(dir, flags)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("\n%s:\n", dir.Name)

	printFiles(dir, files, flags)
	fmt.Print("\n")

	for _, file := range files {
		if file.IsDir() && file.Name != "." && file.Name != ".." {
			nDir, err := getPath(dir.Name + "/" + file.Name)
			if err != nil {
				log.Fatal(err)
			}
			printRecursive(nDir, flags)
		}
	}
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
		fmt.Printf("%s\n", file.Name)
	}
}

func printDirSize(files []FileInfo) {
	var dirSize int

	for _, file := range files {
		if file.Size() < 1024 && file.Size() != 0 {
			dirSize += 1
		} else {
			dirSize += int(math.Ceil(float64(file.Size()) / 1024))
		}
	}

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
		colorLightGreenBackground, file.Name, colorReset)

	return coloredFolder
}
