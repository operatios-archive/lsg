package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/logrusorgru/aurora"
)

func Glob(path string, args Args) {
	fileNames := getFileNames(path)

	parents := make(map[string][]string)
	for _, fileName := range fileNames {
		dir := filepath.Dir(fileName)
		parents[dir] = append(parents[dir], fileName)
	}

	keys := make([]string, 0, len(parents))
	for k := range parents {
		keys = append(keys, k)
	}
	sort.Sort(caseInsensitive(keys))

	for _, parent := range keys {
		if !args.All && isPathHidden(parent) {
			continue
		}

		fmt.Fprintf(BUF_STDOUT, "%s:\n", parent)
		children := getFilesFromStr(parents[parent], args.All)
		processFiles(children, args)
	}
}

func processFiles(files []File, args Args) {
	sortFiles(files, args)

	if args.LongList {
		formatList(files, args)
	} else {
		formatGrid(files, args)
	}
}

func processTree(files []File, fromDepths map[int]bool, args Args) {
	if len(files) == 0 {
		return
	}

	sortFiles(files, args)
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
		fmt.Fprintln(BUF_STDOUT, prefix+file.colored(args))

		if file.isDir() && !file.isLink() {
			subFiles, _ := getFilesFromPath(file.path, args.All)
			processTree(subFiles, fromDepths, args)
		}
	}
}

func getFileNames(pattern string) []string {
	if !strings.Contains(pattern, "**") {
		matches, _ := filepath.Glob(pattern)
		return matches
	}

	matches, _ := doublestar.Glob(pattern)
	return matches
}

// We could just Glob path/* to reduce the amount of code, but this is faster
func getFilesFromPath(path string, showHidden bool) ([]File, error) {
	var result []File

	fileInfos, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		file := fileFromInfo(fileInfo, path)

		if showHidden || !file.isHidden() {
			result = append(result, file)
		}
	}
	return result, nil
}

func getFilesFromStr(fileNames []string, showHidden bool) []File {
	var result []File

	for _, fileName := range fileNames {
		file, err := fileFromStr(fileName)

		if err != nil {
			continue
		}

		if showHidden || !file.isHidden() {
			result = append(result, file)
		}
	}
	return result
}

func isPathHidden(path string) bool {
	components := splitPath(path)

	for i := 1; i <= len(components); i++ {
		f := filepath.Join(components[:i]...)
		if f == "." {
			continue
		}

		if strings.Contains(f, "..") {
			abs, _ := filepath.Abs(f)
			f = abs
		}

		file, err := fileFromStr(f)
		if err != nil {
			continue
		}

		if file.isHidden() {
			return true
		}
	}
	return false
}

func splitPath(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

func humanSize(size int) string {
	if size == 0 {
		return "-"
	}

	if size < 1024 {
		return fmt.Sprintf("%d", size)
	}

	fSize := float64(size)
	fSize /= 1024

	for _, unit := range []rune{'K', 'M', 'G', 'T', 'P', 'E', 'Z'} {
		if fSize < 9 {
			return fmt.Sprintf("%.1f%c", fSize, unit)
		} else if fSize < 1000 {
			return fmt.Sprintf("%.0f%c", fSize, unit)
		}
		fSize /= 1024
	}

	return fmt.Sprintf("%.1f%c", fSize, 'Y')
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
	if rowWidth+(len(columnWidths)-1)*args.ColSep >= TTY_WIDTH {
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

	if args.Columns > 0 {
		columns = args.Columns
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
			sep := strings.Repeat(" ", args.ColSep)
			fmt.Fprintln(BUF_STDOUT, strings.Join(rows[i], sep))
		}
	} else {
		for i := range files {
			fmt.Fprintln(BUF_STDOUT, files[i].colored(args))
		}
	}
}

func formatList(files []File, args Args) {
	var sizes []string
	var sizeTotal int
	var sizeAlign int

	for i := range files {
		var sizeEntry string
		sizeTotal += files[i].size()

		if args.Bytes {
			sizeEntry = fmt.Sprint(files[i].size())
		} else {
			sizeEntry = files[i].sizeHuman()
		}

		if len(sizeEntry) > sizeAlign {
			sizeAlign = len(sizeEntry)
		}
		sizes = append(sizes, sizeEntry)
	}

	if args.Bytes {
		fmt.Fprintf(BUF_STDOUT, "total %d\n", sizeTotal)
	} else {
		fmt.Fprintf(BUF_STDOUT, "total %s\n", humanSize(sizeTotal))
	}

	for i := range files {
		sizeEntry := fmt.Sprintf("%*s", sizeAlign, sizes[i])
		if !args.NoColors {
			sizeEntry = aurora.Colorize(sizeEntry, aurora.GreenFg).String()
		}

		var line string
		if args.FileMode {
			line += files[i].fileMode() + " "
		}
		line += sizeEntry + " "
		line += files[i].modTime() + " "
		line += files[i].colored(args)
		fmt.Fprintln(BUF_STDOUT, line)
	}
}