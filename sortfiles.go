package main

import "sort"

func sortFilesCustom(files []FileInfo) {
	// to avoid . and ..
	startIndex := 0
	if len(files) > 1 && files[0].Name == "." && files[1].Name == ".." {
		startIndex = 2
	} else if len(files) > 0 && (files[0].Name == "." || files[0].Name == "..") {
		startIndex = 1
	}

	sort.Slice(files[startIndex:], func(i, j int) bool {
		i += startIndex
		j += startIndex

		nameI := files[i].Name
		nameJ := files[j].Name

		// trim the non-alphabetical start
		if !isAlphabetical(nameI[0]) {
			nameI = trimNonAlphabeticalPrefix(nameI)
		}
		if !isAlphabetical(nameJ[0]) {
			nameJ = trimNonAlphabeticalPrefix(nameJ)
		}

		return nameI < nameJ
	})
}

// check the char if alpha
func isAlphabetical(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// trim non-alpha
func trimNonAlphabeticalPrefix(s string) string {
	for i := 0; i < len(s); i++ {
		if isAlphabetical(s[i]) {
			return s[i:]
		}
	}
	return s // return original string if no alphabetical character is found
}
