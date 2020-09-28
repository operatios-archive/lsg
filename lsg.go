package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/logrusorgru/aurora"
)

func processGlob(path string, args Args) {
	fileNames := Glob(path)

	parents := make(map[string][]string)
	for _, fileName := range fileNames {
		dir := filepath.Dir(fileName)
		parents[dir] = append(parents[dir], fileName)
	}

	keys := make([]string, 0, len(parents))
	for k := range parents {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return cmpCaseInsensitive(keys[i], keys[j])
	})

	for _, parent := range keys {
		if !args.all && isPathHidden(parent) {
			continue
		}

		children := getParentFiles(parents[parent], args.all)
		if len(children) == 0 {
			continue
		}

		fmt.Fprintf(bufStdout, "%s:\n", parent)
		processFiles(children, args)
	}
}

func processFiles(files []File, args Args) {
	sortFiles(files, args.sort, args.reverse)

	if args.longList {
		formatList(files, args)
	} else {
		formatGrid(files, args)
	}
}

func processTree(files []File, fromDepths map[int]bool, args Args) {
	if len(files) == 0 {
		return
	}

	sortFiles(files, args.sort, args.reverse)
	depth := len(splitPath(files[0].path)) - 1

	for _, file := range files {
		isLast := file == files[len(files)-1]

		if file.isDir() {
			if isLast {
				delete(fromDepths, depth)
			} else {
				fromDepths[depth] = true
			}
		}

		var prefix string
		for i := 0; i < depth; i++ {
			if exists := fromDepths[i]; exists {
				prefix += "│  "
			} else {
				prefix += "   "
			}
		}
		if isLast {
			prefix += "└──"
		} else {
			prefix += "├──"
		}

		fmt.Fprintln(bufStdout, prefix+file.colored(args))

		if file.isDir() && !file.isLink() {
			subFiles, _ := getFiles(file.path, args.all)
			processTree(subFiles, fromDepths, args)
		}
	}
}

func Glob(pattern string) []string {
	var matches []string

	if !strings.Contains(pattern, "**") {
		matches, _ = filepath.Glob(pattern)
	} else {
		matches, _ = doublestar.Glob(pattern)
	}

	return matches
}

func getFiles(path string, showHidden bool) ([]File, error) {
	var result []File

	fileInfos, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		file := File{fileInfo, filepath.Join(path, fileInfo.Name())}

		if showHidden || !file.isHidden() {
			result = append(result, file)
		}
	}
	return result, nil
}

func getParentFiles(fileNames []string, showHidden bool) []File {
	var result []File

	for _, fileName := range fileNames {
		file, err := newFile(fileName)

		if err != nil {
			continue
		}

		if showHidden || !file.isHidden() {
			result = append(result, file)
		}
	}
	return result
}

func getRowCol(i int, rows int) (int, int) {
	row := i % rows
	return row, (i - row) / rows
}

func formatRows(files []File, columns int, args Args) [][]string {
	var rows int
	if len(files)%columns != 0 {
		rows = (len(files) / columns) + 1
	} else {
		rows = len(files) / columns
	}

	rowSlice := make([][]string, rows)
	columnWidths := make([]int, columns)

	for i, file := range files {
		_, col := getRowCol(i, rows)

		nameLength := utf8.RuneCountInString(file.pretty(args))
		if nameLength > columnWidths[col] {
			columnWidths[col] = nameLength
		}
	}

	var rowWidth int
	for _, width := range columnWidths {
		rowWidth += width
	}
	if rowWidth+(len(columnWidths)-1)*args.colSep >= terminalWidth {
		return nil
	}

	for i, file := range files {
		row, col := getRowCol(i, rows)
		wsAmt := columnWidths[col] - utf8.RuneCountInString(file.pretty(args))
		padding := strings.Repeat(" ", wsAmt)

		rowSlice[row] = append(rowSlice[row], file.colored(args)+padding)
	}
	return rowSlice
}

func formatGrid(files []File, args Args) {
	columns := 2
	goingBackwards := false

	if args.columns > 0 {
		columns = args.columns
		goingBackwards = true
	}

	var rows [][]string
	for columns > 1 {
		rows = formatRows(files, columns, args)
		if goingBackwards && rows != nil {
			break
		}

		if rows == nil || columns > len(files) {
			goingBackwards = true
		}

		if !goingBackwards {
			columns *= 2
		} else {
			columns--
		}
	}

	if columns > 1 {
		for i := range rows {
			sep := strings.Repeat(" ", args.colSep)
			fmt.Fprintln(bufStdout, strings.Join(rows[i], sep))
		}
	} else {
		for i := range files {
			fmt.Fprintln(bufStdout, files[i].colored(args))
		}
	}
}

func formatList(files []File, args Args) {
	var sizes []string
	var totalSize int64

	var align struct {
		size     int
		fileMode int
		nLink    int
		owner    int
		group    int
	}

	for _, file := range files {
		var sizeEntry string
		totalSize += file.size()

		if args.bytes {
			sizeEntry = fmt.Sprint(file.size())
		} else {
			sizeEntry = humanizeSize(file.size())
		}
		sizes = append(sizes, sizeEntry)

		if len(sizeEntry) > align.size {
			align.size = len(sizeEntry)
		}

		if args.listExtend {
			if len(file.fileMode()) > align.fileMode {
				align.fileMode = len(file.fileMode())
			}

			nLinkLen := len(fmt.Sprint(file.nLink()))
			if nLinkLen > align.nLink {
				align.nLink = nLinkLen
			}
		}

		if args.listExtend && runtime.GOOS == "linux" {
			if len(file.owner()) > align.owner {
				align.owner = len(file.owner())
			}
			if len(file.group()) > align.group {
				align.group = len(file.group())
			}
		}
	}

	if args.bytes {
		fmt.Fprintf(bufStdout, "total %d\n", totalSize)
	} else {
		fmt.Fprintf(bufStdout, "total %s\n", humanizeSize(totalSize))
	}

	for i, file := range files {
		var line string
		if args.listExtend {
			line += fmt.Sprintf("%-*s ", align.fileMode, file.fileMode())
			line += fmt.Sprintf("%*d ", align.nLink, file.nLink())
		}

		if args.listExtend && runtime.GOOS == "linux" {
			owner := file.owner()
			group := file.group()

			// WSL: file owner of /mnt/ is ""
			if owner == "" {
				owner = group
			}

			line += fmt.Sprintf("%*s ", align.owner-1, owner)
			line += fmt.Sprintf("%*s ", align.group-1, group)
		}

		sizeEntry := fmt.Sprintf("%*s", align.size, sizes[i])
		if !args.noColors {
			sizeEntry = aurora.Colorize(sizeEntry, aurora.GreenFg).String()
		}
		line += sizeEntry + " "
		line += files[i].modTime() + " "
		line += files[i].colored(args)

		fmt.Fprintln(bufStdout, line)
	}
}
