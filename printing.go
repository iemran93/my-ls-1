package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

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
