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
	// Pads the timestamps, which is date + time
	TimestampPadding = iota
	// Pads the date only
	DatePadding
	// Pads the Levels
	LevelPadding
)

// Column stuff

// Represents a column in a logged message.
type Column func(context Context) string

// Logger stuff

const (
	// Successful log, should be universal
	Success = iota
	// The level passed isn't in this Logger
	InvalidLevel
	// There were no Column s set in the Logger
	NoColumnsSet
)

// Represents a Logger that can log to a variety of things.
type Logger interface {
	// Returns all the Levels for this Logger.
	GetLevels() map[string]func() string
	// Returns all the Padding this Logger uses.
	GetPaddings() []Padding
	// Returns all the Column s this Logger uses.
	GetColumns() []Column
	// Logs a message. This should return a status code and an optional error message.
	Log(level, message string) (int, error)
	// Similar to Logger.Log, but can contain additional information if the logger needs it.
	LogWithExtraInfo(level, message string, info interface{}) (int, error)
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

// Represents the Context of a Logger message. This contains the Message being sent,
// the Time of the message, it's Level, and the corresponding Logger.
type Context struct {
	Message string
	Time    time.Time
	Level   string
	Logger  Logger
}

// ANSI color codes

var ansi = regexp.MustCompile("\\x1B(?:[@-Z\\\\-_]|\\[[0-?]*[ -/]*[@-~])")

/*
	Formats the Level by it's display and Padding. If you use the default settings
	then a message would take the shape of:

		Saturday August 29, 2020 @ 5:41:00 | INFO | Hello, world

	This method can add optional padding so that all Levels will line up.
	Without padding:

		Saturday August 29, 2020 @ 5:41:00 | WARNING | Hello, world
		Saturday August 29, 2020 @ 5:41:20 | INFO | Hello, world

	With padding:

		Saturday August 29, 2020 @ 5:41:00 | WARNING | Hello, world
		Saturday August 29, 2020 @ 5:41:20 | INFO    | Hello, world
*/
func (c *Context) FormatLevel() string {
	padding := findPadding(c.Logger.GetPaddings(), LevelPadding) != -1
	after := ""
	display := c.Logger.GetLevels()[c.Level]()
	if padding {
		longest := 0
		for _, element := range c.Logger.GetLevels() {
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

// This formats the date so that it's always as long as the longest date you can display (without repeating).
// The longest date I was able to find is "Wednesday September 30th, 9999"
func (c *Context) FormatDate(layout string) string {
	padding := findPadding(c.Logger.GetPaddings(), DatePadding) != -1
	after := ""
	formatted := c.Time.Format(layout)

	if padding {
		longest := 30 // "Wednesday September 30th, 9999" was the longest date I could find.
		for i := len(formatted); i < longest; i++ {
			after += " "
		}
	}

	return fmt.Sprintf("%v%v", formatted, after)
}

// Formats the time with the layout provided.
func (c *Context) FormatTime(layout string) string {
	return c.Time.Format(layout)
}

var longestTimestampSeen = 0

/*
	This formats the timestamp with the layout provided.
	This padding can change due to the fact that timestamps aren't universal in how they're formatted.
*/
func (c *Context) FormatTimestamp(layout string) string {
	padding := findPadding(c.Logger.GetPaddings(), TimestampPadding) != -1
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
