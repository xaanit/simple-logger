/*
 * Copyright 2020 Jacob Frazier
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
 * associated documentation files (the "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
 * of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following
 * conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial
 * portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
 * PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 * LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT
 * OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */
package simple_logger

import (
	"fmt"
	"regexp"
	"time"
)

// Padding
type Padding int

const (
	TimestampPadding = iota
	DatePadding
	LevelPadding
)

// Column stuff

type Column func(context Context) string

// Logger stuff

type Logger interface {
	GetLevels() map[string]func() string
	GetPaddings() []Padding
	GetColumns() []Column
	Log(level, message string)
}

// Context stuff

const ( // date formatting
	Weekday = "Monday"
	Day     = "2"
	Month   = "January"
	Year    = "2006"
)

const ( // Time formatting
	Hour   = "3"
	Hour24 = "15"
	Minute = "04"
	Second = "05"
)

type Context struct {
	Message string
	Time    time.Time
	Level   string
	logger  Logger
}

// ANSI color codes

var ansi = regexp.MustCompile("\\x1B(?:[@-Z\\\\-_]|\\[[0-?]*[ -/]*[@-~])")

func (c *Context) FormatLevel() string {
	padding := findPadding(c.logger.GetPaddings(), LevelPadding) != -1
	after := ""
	display := c.logger.GetLevels()[c.Level]()
	if padding {
		longest := 0
		for _, element := range c.logger.GetLevels() {
			if l := len(ansi.ReplaceAllString(element(), "")); l > longest {
				longest = l
			}
		}
		for i := len(ansi.ReplaceAllString(display, "")); i < longest; i++ {
			after += " "
		}
	}

	return fmt.Sprintf("%v%v", display, after)
}

func (c *Context) FormatDate(layout string) string {
	padding := findPadding(c.logger.GetPaddings(), DatePadding) != -1
	after := ""
	formatted := c.Time.Format(layout)

	if padding {
		longest := 30 // "Wednesday September 30th, 9999" was the longest date I could findPadding.
		for i := len(formatted); i < longest; i++ {
			after += " "
		}
	}

	return fmt.Sprintf("%v%v", formatted, after)
}

func (c *Context) FormatTime(layout string) string {
	return c.Time.Format(layout)
}

var longestTimestampSeen = 0

func (c *Context) FormatTimestamp(layout string) string {
	padding := findPadding(c.logger.GetPaddings(), TimestampPadding) != -1
	formatted := c.Time.Format(layout)
	after := ""

	if padding {
		i := len(formatted)
		if i > longestTimestampSeen {
			longestTimestampSeen = i
		}

		for ; i < longestTimestampSeen; i++ {
			after += " "
		}
	}
	return fmt.Sprintf("%v%v", formatted, after)
}
