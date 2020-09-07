package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func sortFiles(files []File, args Args) {
	switch strings.ToLower(args.Sort) {
	case "s", "size":
		sort.Sort(sort.Reverse(bySize(files)))
	case "t", "time":
		sort.Sort(sort.Reverse(byTime(files)))
	case "x", "extension":
		sort.Sort(byExtension(files))
	case "c", "category":
		sort.Sort(byCategory(files))
	case "":
		sort.Sort(byName(files))
	default:
		fmt.Fprintf(os.Stderr, "Invalid sorting parameter: %s\n", args.Sort)
		os.Exit(1)
	}

	if args.Reverse {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}
}

func caseInsensitiveSort(a, b string) bool {
	return strings.ToLower(a) < strings.ToLower(b)
}

type caseInsensitive []string

func (f caseInsensitive) Len() int {
	return len(f)
}
func (f caseInsensitive) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f caseInsensitive) Less(i, j int) bool {
	return caseInsensitiveSort(f[i], f[j])
}

type byName []File

func (f byName) Len() int {
	return len(f)
}
func (f byName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byName) Less(i, j int) bool {
	return caseInsensitiveSort(f[i].name(), f[j].name())
}

type byCategory []File

func (f byCategory) Len() int {
	return len(f)
}
func (f byCategory) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byCategory) Less(i, j int) bool {
	return f[i].category() < f[j].category()
}

type byExtension []File

func (f byExtension) Len() int {
	return len(f)
}
func (f byExtension) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byExtension) Less(i, j int) bool {
	return caseInsensitiveSort(f[i].ext(), f[j].ext())
}

type bySize []File

func (f bySize) Len() int {
	return len(f)
}
func (f bySize) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f bySize) Less(i, j int) bool {
	return f[i].size() < f[j].size()
}

type byTime []File

func (f byTime) Len() int {
	return len(f)
}
func (f byTime) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byTime) Less(i, j int) bool {
	return f[i].modTime() < f[j].modTime()
}
