package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

const (
	helpAll       = "do not ignore hidden files"
	helpLongList  = "use a long listing format"
	helpBytes     = "with -l: print size in bytes"
	helpExtend    = "with -l: print filemode and owner/group info"
	helpTree      = "use a tree format"
	helpSort      = "sort by size (s), time (t), extension (x), category (c)"
	helpReverse   = "reverse file order"
	helpColumns   = "set maximum amount of columns"
	helpColSep    = "set column separator length"
	helpNoTargets = "disable link targets"
	helpNoColors  = "disable colors"
	helpNoIcons   = "disable icons"
	helpShow      = "show this message and exit"
)

type Args struct {
	paths      []string
	all        bool
	longList   bool
	bytes      bool
	listExtend bool
	tree       bool
	sort       string
	reverse    bool
	columns    int
	colSep     int
	noTargets  bool
	noColors   bool
	noIcons    bool
}

func getArgs() Args {
	args := Args{}

	flag.CommandLine.SortFlags = false

	flag.BoolVarP(&args.all, "all", "a", false, helpAll)
	flag.BoolVarP(&args.longList, "long-listing", "l", false, helpLongList)
	flag.BoolVarP(&args.bytes, "bytes", "b", false, helpBytes)
	flag.BoolVarP(&args.listExtend, "extend", "x", false, helpExtend)
	flag.BoolVarP(&args.tree, "tree", "t", false, helpTree)
	flag.StringVarP(&args.sort, "sort", "s", "", helpSort)
	flag.BoolVarP(&args.reverse, "reverse", "r", false, helpReverse)
	flag.IntVarP(&args.columns, "columns", "c", 0, helpColumns)
	flag.IntVar(&args.colSep, "col-sep", 2, helpColSep)
	flag.BoolVar(&args.noTargets, "no-targets", false, helpNoTargets)
	flag.BoolVar(&args.noColors, "no-colors", false, helpNoColors)
	flag.BoolVar(&args.noIcons, "no-icons", false, helpNoIcons)

	var showHelp *bool = flag.BoolP("help", "h", false, helpShow)

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if args.colSep < 0 {
		fmt.Fprintln(os.Stderr, "column separator length should be >=0")
		os.Exit(1)
	}

	args.paths = flag.Args()
	return args
}
