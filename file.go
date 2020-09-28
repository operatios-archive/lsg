package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/operatios/lsg/category"
	"github.com/operatios/lsg/icons"
)

type File struct {
	info os.FileInfo
	path string
}

func newFile(path string) (File, error) {
	fileInfo, err := os.Lstat(path)

	// File got deleted while executing
	if err != nil {
		return File{}, err
	}

	return File{fileInfo, path}, nil
}

func (f File) name() string {
	return f.info.Name()
}

func (f File) ext() string {
	return filepath.Ext(f.name())
}

func (f File) size() int64 {
	return f.info.Size()
}

func (f File) modTime() string {
	modtime := f.info.ModTime()
	if modtime.Year() == time.Now().Year() {
		return modtime.Format("Jan 02 15:04")
	}
	return modtime.Format("Jan 02  2006")
}

func (f File) fileMode() string {
	return f.info.Mode().String()
}

func (f File) isLink() bool {
	return f.info.Mode()&os.ModeSymlink != 0
}

func (f File) isBroken() bool {
	target, _ := filepath.EvalSymlinks(f.path)
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

	if !args.noTargets && f.isLink() {
		var arrow string
		if args.noIcons {
			arrow = "->"
		} else {
			arrow = icons.LinkArrow
		}
		displayName = displayName + " " + arrow + " " + f.target()
	}

	if !args.noIcons {
		displayName = f.icon() + " " + displayName
	}

	return displayName
}

func (f File) category() int {
	if f.isLink() {
		if f.isBroken() {
			return category.Broken
		}
		return category.Symlink
	}

	if f.isDir() {
		return category.Dir
	}

	if extCategory, ok := category.Extensions[f.ext()]; ok {
		return extCategory
	}

	return category.File
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

	if icon, ok := icons.Extensions[f.ext()]; ok {
		return icon
	}

	return icons.File
}

func (f File) colored(args Args) string {
	pretty := f.pretty(args)

	if args.noColors {
		return pretty
	}
	return aurora.Colorize(pretty, colorScheme[f.category()]).String()
}
