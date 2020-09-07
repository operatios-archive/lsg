package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

const (
	HELP_ALL       = "do not ignore hidden files"
	HELP_LONGLIST  = "use a long listing format"
	HELP_BYTES     = "with -l: print size in bytes"
	HELP_FILEMODE  = "with -l: print filemode"
	HELP_TREE      = "use a tree format"
	HELP_SORT      = "sort by [size|s|time|t|extension|x|category|c]"
	HELP_REVERSE   = "reverse file order"
	HELP_COLUMNS   = "set maximum amount of columns"
	HELP_COLSEP    = "set column separator length"
	HELP_NOTARGETS = "disable link targets"
	HELP_NOCOLORS  = "disable colors"
	HELP_NOICONS   = "disable icons"
	HELP_SHOW      = "show this message and exit"
)

type Args struct {
	Paths     []string
	All       bool
	LongList  bool
	Bytes     bool
	FileMode  bool
	Tree      bool
	Sort      string
	Reverse   bool
	Columns   int
	ColSep    int
	NoTargets bool
	NoColors  bool
	NoIcons   bool
}

func getArgs() Args {
	args := Args{}
	flag.CommandLine.SortFlags = false

	flag.BoolVarP(&args.All,      "all",          "a",   false, HELP_ALL)
	flag.BoolVarP(&args.LongList, "long-listing", "l",   false, HELP_LONGLIST)
	flag.BoolVarP(&args.Bytes,    "bytes",        "b",   false, HELP_BYTES)
	flag.BoolVarP(&args.FileMode, "filemode",     "f",   false, HELP_FILEMODE)
	flag.BoolVarP(&args.Tree,     "tree",         "t",   false, HELP_TREE)
	flag.StringVarP(&args.Sort,   "sort",         "s",   "",    HELP_SORT)
	flag.BoolVarP(&args.Reverse,  "reverse",      "r",   false, HELP_REVERSE)
	flag.IntVarP(&args.Columns,   "columns",      "c",   0,     HELP_COLUMNS)
	flag.IntVar(&args.ColSep,     "col-sep",             2,     HELP_COLSEP)
	flag.BoolVar(&args.NoTargets, "no-targets",          false, HELP_NOTARGETS)
	flag.BoolVar(&args.NoColors,  "no-colors",           false, HELP_NOCOLORS)
	flag.BoolVar(&args.NoIcons,   "no-icons",            false, HELP_NOICONS)

	var showHelp *bool = flag.BoolP("help", "h", false, HELP_SHOW)
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if args.ColSep < 0 {
		fmt.Fprintln(os.Stderr, "column separator length should be >=0")
		os.Exit(1)
	}

	args.Paths = flag.Args()
	if len(args.Paths) == 0 {
		args.Paths = append(args.Paths, ".")
	}

	return args
}
