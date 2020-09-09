package main

import (
	"os"
	"syscall"
)

var (
	msvcrt         = syscall.NewLazyDLL("msvcrt.dll")
	_get_osfhandle = msvcrt.NewProc("_get_osfhandle")
)

func (f File) attrs() uint32 {
	return f.fileInfo.Sys().(*syscall.Win32FileAttributeData).FileAttributes
}

func (f File) isDir() bool {
	if f.isLink() {
		return f.attrs()&syscall.FILE_ATTRIBUTE_DIRECTORY != 0
	}
	return f.fileInfo.IsDir()
}

func (f File) isHidden() bool {
	dotHidden := f.name()[0] == '.'
	if dotHidden {
		return true
	}
	return f.attrs()&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

func (f File) nLink() int {
	file, err := os.Open(f.path)
	if err != nil {
		return 1
	}

	_, handle, _ := _get_osfhandle.Call(file.Fd())
	var info syscall.ByHandleFileInformation
	syscall.GetFileInformationByHandle(syscall.Handle(handle), &info)

	if info.NumberOfLinks > 0 {
		return int(info.NumberOfLinks)
	}
	return 1
}

func (f File) owner() string {
	return ""
}

func (f File) group() string {
	return ""
}
