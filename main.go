package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// To increase Stdout print speed we do buffered output
	// So instead of fmt.Println we use fmt.Fprintln(BUF_STDOUT, ...)
	// BUT, if an error causes program to os.Exit, use fmt.Println
	BUF_STDOUT *bufio.Writer = bufio.NewWriter(os.Stdout)

	ISATTY          = terminal.IsTerminal(int(os.Stdout.Fd()))
	TTY_WIDTH, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
)

func main() {
	defer BUF_STDOUT.Flush()
	args := getArgs()

	if !ISATTY {
		args.Columns = 1
		args.NoColors = true
		args.NoIcons = true
	}

	if runtime.GOOS == "windows" && !args.NoColors {
		enableColors()
	}

	if args.Tree {
		doTree(args)
	} else {
		doLS(args)
	}
}

func doTree(args Args) {
	wd, _ := os.Getwd()

	for _, path := range args.Paths {

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

		files, _ := getFilesFromPath(".", args.All)

		clean := filepath.Clean(path)
		if !args.NoColors {
			clean = aurora.Colorize(clean, colorScheme[DIR]).String()
		}
		fmt.Fprintln(BUF_STDOUT, clean)

		processTree(files, map[int]bool{0: true}, args)
	}
}

func doLS(args Args) {
	for _, path := range args.Paths {
		if strings.ContainsRune(path, '*') {
			Glob(path, args)

		} else {
			files, err := getFilesFromPath(path, args.All)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			if len(args.Paths) > 1 {
				fmt.Fprintln(BUF_STDOUT, filepath.Clean(path)+":")
			}
			processFiles(files, args)
		}
	}
}
