package main

import "syscall"

func (f File) attrs() uint32 {
	pUTF16, _ := syscall.UTF16PtrFromString(f.path)
	attrs, _ := syscall.GetFileAttributes(pUTF16)
	return attrs
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
