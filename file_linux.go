package main

import (
	"fmt"
	"os/user"
	"syscall"
)

func (f File) isDir() bool {
	return f.fileInfo.IsDir()
}

func (f File) isHidden() bool {
	return f.name()[0] == '.'
}

func (f File) stat_t() syscall.Stat_t {
	return *f.fileInfo.Sys().(*syscall.Stat_t)
}

func (f File) group() string {
	group, err := user.LookupGroupId(fmt.Sprint(f.stat_t().Gid))
	if err != nil {
		return ""
	}
	return group.Name
}

func (f File) owner() string {
	user, err := user.LookupId(fmt.Sprint(f.stat_t().Uid))
	if err != nil {
		return ""
	}
	return user.Name
}

func (f File) nLink() int {
	return int(f.stat_t().Nlink)
}
