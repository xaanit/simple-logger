package simple_logger

import (
	"github.com/logrusorgru/aurora/v3"
)

var (
	Info    = func() string { return aurora.Cyan("INFO").String() }
	Warning = func() string { return aurora.Yellow("WARNING").String() }
	Debug   = func() string { return aurora.Green("DEBUG").String() }
	Error   = func() string { return aurora.Red("ERROR").String() }
	Fatal   = func() string { return aurora.Bold(aurora.Red("FATAL")).String() }
)
