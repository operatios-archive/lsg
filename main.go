package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/operatios/lsg/category"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	bufStdout           = bufio.NewWriter(os.Stdout)
	terminalWidth, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
)

func isatty() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func main() {
	defer bufStdout.Flush()

	args := getArgs()

	if len(args.paths) == 0 {
		args.paths = append(args.paths, ".")
	}

	if !isatty() {
		args.columns = 1
		args.noColors = true
		args.noIcons = true
	}

	if runtime.GOOS == "windows" && !args.noColors {
		err := enableColors()
		if err != nil {
			log.Fatal(err)
		}
	}

	if args.tree {
		doTree(args)
	} else {
		doLS(args)
	}
}

func doLS(args Args) {
	for _, path := range args.paths {

		if strings.ContainsRune(path, '*') {
			processGlob(path, args)
		} else {
			files, err := getFiles(path, args.all)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			if len(args.paths) > 1 {
				fmt.Fprintln(bufStdout, filepath.Clean(path)+":")
			}

			processFiles(files, args)
		}
	}
}

func doTree(args Args) {
	wd, _ := os.Getwd()

	for _, path := range args.paths {

		var err error
		if filepath.IsAbs(path) {
			err = os.Chdir(path)
		} else {
			err = os.Chdir(filepath.Join(wd, path))
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		files, _ := getFiles(".", args.all)

		clean := filepath.Clean(path)
		if !args.noColors {
			clean = aurora.Colorize(clean, colorScheme[category.Dir]).String()
		}
		fmt.Fprintln(bufStdout, clean)

		processTree(files, map[int]bool{0: true}, args)
	}
}
