package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/operatios/lsg/icons"
)

type File struct {
	fileInfo os.FileInfo
	path     string
}

func fileFromStr(path string) (File, error) {
	fileInfo, err := os.Lstat(path)
	// File got deleted while executing
	if err != nil {
		return File{}, err
	}

	return File{fileInfo, path}, nil
}

func fileFromInfo(info os.FileInfo, path string) File {
	return File{info, filepath.Join(path, info.Name())}
}

func (f File) name() string {
	return f.fileInfo.Name()
}

func (f File) ext() string {
	return filepath.Ext(f.name())
}

func (f File) size() int {
	return int(f.fileInfo.Size())
}

func (f File) sizeHuman() string {
	return humanSize(f.size())
}

func (f File) modTime() string {
	return f.fileInfo.ModTime().Format("Jan 02 15:04")
}

func (f File) fileMode() string {
	return f.fileInfo.Mode().String()
}

func (f File) isLink() bool {
	return f.fileInfo.Mode()&os.ModeSymlink != 0
}

func (f File) isBroken() bool {
	target, _ := os.Readlink(f.path)
	_, err := os.Stat(target)
	return err != nil
}

func (f File) target() string {
	target, _ := os.Readlink(f.path)
	wd, _ := os.Getwd()
	relPath, _ := filepath.Rel(wd, target)

	if relPath != "" && !strings.HasPrefix(relPath, "..") {
		return relPath
	}
	return target
}

func (f File) pretty(args Args) string {
	displayName := f.name()

	if !args.NoTargets && f.isLink() {
		var arrow string
		if args.NoIcons {
			arrow = "->"
		} else {
			arrow = icons.LinkArrow
		}
		displayName = displayName + " " + arrow + " " + f.target()
	}

	if !args.NoIcons {
		displayName = f.icon() + " " + displayName
	}

	return displayName
}

func (f File) category() int {
	if f.isLink() {
		if f.isBroken() {
			return BROKEN
		}
		return SYMLINK
	}

	if f.isDir() {
		return DIR
	}

	if extensionCategory, ok := extensionToCategory[f.ext()]; ok {
		return extensionCategory
	}

	return FILE
}

func (f File) icon() string {
	if f.isLink() {
		if f.isDir() {
			return icons.LinkDir
		}
		return icons.LinkFile
	}

	if f.isDir() {
		return icons.Dir
	}

	icon := icons.Extensions[f.ext()]
	if icon != "" {
		return icon
	}

	return icons.File
}

func (f File) colored(args Args) string {
	if args.NoColors {
		return f.pretty(args)
	}

	color := colorScheme[f.category()]
	return aurora.Colorize(f.pretty(args), color).String()
}
