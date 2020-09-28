package main

import (
	"github.com/logrusorgru/aurora"
	"github.com/operatios/lsg/category"
)

var colorScheme = map[int]aurora.Color{
	category.File:       aurora.GreenFg,
	category.Dir:        aurora.BlueFg + aurora.BoldFm,
	category.Symlink:    aurora.CyanFg + aurora.UnderlineFm,
	category.Broken:     aurora.BlackFg + aurora.RedBg,
	category.Archive:    aurora.RedFg + aurora.UnderlineFm,
	category.Executable: aurora.RedFg + aurora.BoldFm,
	category.Code:       aurora.MagentaFg,
	category.Image:      aurora.YellowFg,
	category.Audio:      aurora.YellowFg + aurora.BoldFm,
	category.Video:      aurora.YellowFg,
}
