package main

import (
	"fmt"
	"log"
	"os/user"
	"syscall"
)

func (f File) isDir() bool {
	return f.info.IsDir()
}

func (f File) isHidden() bool {
	return f.name()[0] == '.'
}

func (f File) stat_t() syscall.Stat_t {
	return *f.info.Sys().(*syscall.Stat_t)
}

func (f File) group() string {
	group, err := user.LookupGroupId(fmt.Sprint(f.stat_t().Gid))
	if err != nil {
		log.Panic(err)
	}
	return group.Name
}

func (f File) owner() string {
	user, err := user.LookupId(fmt.Sprint(f.stat_t().Uid))
	if err != nil {
		log.Panic(err)
	}
	return user.Name
}

func (f File) nLink() uint {
	return uint(f.stat_t().Nlink)
}
