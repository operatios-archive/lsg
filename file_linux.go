package main

func (f File) isDir() bool {
	return f.fileInfo.IsDir()
}

func (f File) isHidden() bool {
	return f.name()[0] == '.'
}
