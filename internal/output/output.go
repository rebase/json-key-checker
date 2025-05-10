package output

import (
	"github.com/fatih/color"
)

var (
	SuccessColor    = color.New(color.FgGreen).SprintFunc()
	WarningColor    = color.New(color.FgYellow).SprintFunc()
	ErrorColor      = color.New(color.FgRed).SprintFunc()
	InfoColor       = color.New(color.FgCyan).SprintFunc()
	Bold            = color.New(color.Bold).SprintFunc()
	MissingKeyColor = color.New(color.FgRed).SprintFunc()
	HeaderColor     = color.New(color.FgBlue).SprintFunc()
)
