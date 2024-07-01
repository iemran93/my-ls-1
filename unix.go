package main

import (
	"fmt"
	"os/user"
	"syscall"
)

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
