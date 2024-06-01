//go:build unix
// +build unix

package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
	"time"
)

func printLongFormat(file os.FileInfo) {
	stat, ok := file.Sys().(*syscall.Stat_t)
	if !ok {
		fmt.Printf("Error: unable to retrieve file info for %s\n", file.Name())
		return
	}

	userInfo, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		userInfo = &user.User{Username: "unknown"}
	}

	groupInfo, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
	if err != nil {
		groupInfo = &user.Group{Name: "unknown"}
	}

	modTime := file.ModTime().Format(time.RFC3339)
	fmt.Printf("%s\t%d\t%s\t%s\t%d\t%s\t%s\n",
		file.Mode().String(),
		stat.Nlink,
		userInfo.Username,
		groupInfo.Name,
		file.Size(),
		modTime,
		file.Name())
}
