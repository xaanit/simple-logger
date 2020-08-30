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
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
)

// Defines the methods needed to build a Logger instance
type LoggerBuilder interface {
	// Adds a level for the Logger to use. The name is called from the Logger.Log function
	// and the display is used in the actual logging.
	AddLevel(name string, display func() string) LoggerBuilder
	// Adds a level of padding to the Logger. This only works with default Padding.
	AddPadding(padding Padding) LoggerBuilder
	// Adds a new Column to the logger. This should always add to the end of the list.
	AddColumn(column Column) LoggerBuilder
	// Adds a new Column into the index passed. This should error if the index does not exist,
	// but should add to the if the length of the underlying array is passed.
	AddColumnByIndex(index uint, column Column) (LoggerBuilder, error)
	// Builds a new Logger instance.
	Build() Logger
}

/* Sets the default for this builder. The defaults are outlined below:

 	Levels:
	 +---------+----------+--------------------------+
	 | Name    | Variable | Colours                  |
	 +---------+----------+--------------------------+
	 | INFO    | Info     | aurora.Cyan              |
	 +---------+----------+--------------------------+
	 | DEBUG   | Debug    | aurora.Green             |
	 +---------+----------+--------------------------+
	 | ERROR   | Error    | aurora.Red               |
	 +---------+----------+--------------------------+
	 | FATAL   | Fatal    | aurora.Bold + aurora.Red |
	 +---------+----------+--------------------------+
	 | WARNING | Warning  | aurora.Yellow            |
	 +---------+----------+--------------------------+

 	Paddings:
 		- TimestampPadding
 		- LevelPadding

 	Columns:
 		+-----------+-------+---------+
 		| Timestamp | Level | Message |
 		+-----------+-------+---------+


 	The last three arguments are for excluding certain defaults.
 	If you'd like to only have the Info, Warning, and Error levels you'd pass []string{"FATAL", "DEBUG"}.

 	Columns in this way are indexed based starting from 0. If you'd like to remove the timestamp portion you'd pass []int{0}
*/
func SetDefaults(builder LoggerBuilder, excludeLevels []string, excludePaddings []Padding, excludeColumns []uint) LoggerBuilder {
	defaultLevels := func(level string, display func() string) {
		if findStrings(excludeLevels, level) == -1 {
			builder.AddLevel(level, display)
		}
	}

	defaultLevels("INFO", Info)
	defaultLevels("DEBUG", Debug)
	defaultLevels("ERROR", Error)
	defaultLevels("FATAL", Fatal)
	defaultLevels("WARNING", Warning)

	defaultPaddings := func(padding Padding) {
		if findPadding(excludePaddings, padding) == -1 {
			builder.AddPadding(padding)
		}
	}

	defaultPaddings(TimestampPadding)
	defaultPaddings(LevelPadding)

	defaultColumns := func(column uint, display Column) {
		if findInts(excludeColumns, column) == -1 {
			_, _ = builder.AddColumnByIndex(column, display)
		}
	}

	defaultColumns(0, func(context Context) string {
		layout := fmt.Sprintf("%v %v %v, %v @ %v:%v:%v", Weekday, Month, Day, Year, Hour, Minute, Second)
		return aurora.BrightBlue(context.FormatTimestamp(layout)).String()
	})
	defaultColumns(1, func(context Context) string { return context.FormatLevel() })
	defaultColumns(2, func(context Context) string { return context.Message })

	return builder
}

// Makes a new GenericLoggerBuilder to be used for most loggers. Since golang doesn't allow for
// default methods in interfaces.
func NewGenericLoggerBuilder() *GenericLoggerBuilder {
	return &GenericLoggerBuilder{
		levels:   make(map[string]func() string),
		paddings: make(map[Padding]interface{}),
		columns:  make([]Column, 0),
	}
}

type GenericLoggerBuilder struct {
	levels   map[string]func() string
	paddings map[Padding]interface{}
	columns  []Column
}

// Implements LoggerBuilder.AddLevel
func (b *GenericLoggerBuilder) AddLevel(name string, display func() string) {
	b.levels[name] = display
}

// Implements LoggerBuilder.AddPadding
func (b *GenericLoggerBuilder) AddPadding(padding Padding) {
	b.paddings[padding] = nil
}

// Implements LoggerBuilder.AddColumn
func (b *GenericLoggerBuilder) AddColumn(column Column) {
	b.columns = append(b.columns, column)
}

// Implements LoggerBuilder.AddColumnByIndex
func (b *GenericLoggerBuilder) AddColumnByIndex(index uint, column Column) error {
	if index == uint(len(b.columns)) {
		b.AddColumn(column)
		return nil
	}
	if index > uint(len(b.columns)) {
		return errors.New("index must be less than the index of the array")
	}

	// https://stackoverflow.com/questions/46128016/insert-a-value-in-a-slice-at-a-given-index
	b.columns = append(b.columns, nil)
	copy(b.columns[index+1:], b.columns[index:])
	b.columns[index] = column
	return nil
}
