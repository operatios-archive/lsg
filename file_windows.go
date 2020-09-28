package main

import (
	"log"
	"syscall"
)

func (f File) attrs() uint32 {
	return f.info.Sys().(*syscall.Win32FileAttributeData).FileAttributes
}

func (f File) isDir() bool {
	if f.isLink() {
		return f.attrs()&syscall.FILE_ATTRIBUTE_DIRECTORY != 0
	}
	return f.info.IsDir()
}

func (f File) isHidden() bool {
	if f.name()[0] == '.' {
		return true
	}
	return f.attrs()&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

func (f File) nLink() uint {
	h, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(f.path),
		0,
		0,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_OPEN_REPARSE_POINT|syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0)
	defer syscall.Close(h)

	if err != nil {
		log.Panic(err)
	}

	var info syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(h, &info)

	if err != nil {
		log.Panic(err)
	}
	return uint(info.NumberOfLinks)
}

func (f File) owner() string {
	return ""
}

func (f File) group() string {
	return ""
}
