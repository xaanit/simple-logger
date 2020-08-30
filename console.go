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
	"time"
)

// A Logger implementation that logs to console using fmt.Println
type ConsoleLogger struct {
	levels   map[string]func() string
	paddings []Padding
	columns  []Column
}

func (c ConsoleLogger) createContext(level, message string) Context {
	return Context{
		Message: message,
		Time:    time.Now(),
		Level:   level,
		logger:  c,
	}
}

// Implements Logger.GetLevels
func (c ConsoleLogger) GetLevels() map[string]func() string {
	return c.levels
}

// Implements Logger.GetPaddings
func (c ConsoleLogger) GetPaddings() []Padding {
	return c.paddings
}

// Implements Logger.GetColumns
func (c ConsoleLogger) GetColumns() []Column {
	return c.columns
}

/*
	Implements Logger.Log.

	Returns
		- Success when there was no problems
		- InvalidLevel when the level provided isn't in this Logger
		- NoColumnsSet when there are no columns set for this Logger
*/
func (c ConsoleLogger) Log(level, message string) (int, error) {
	if _, ok := c.GetLevels()[level]; !ok {
		return InvalidLevel, errors.New(fmt.Sprintf("%v is not a valid level for this logger", level))
	}

	if len(c.GetColumns()) == 0 {
		return NoColumnsSet, errors.New("you must set at least one column")
	}

	context := c.createContext(level, message)
	format := ""
	for index, column := range c.columns {
		pad := ""
		if index != 0 {
			pad = " "
		}

		end := ""
		if index != len(c.columns)-1 {
			end = " |"
		}

		format += fmt.Sprintf("%v%v%v", pad, column(context), end)
	}

	fmt.Println(format)

	return Success, nil
}

// Creates a new LoggerBuilder for making instances of ConsoleLogger
func ConsoleLoggerBuilder() LoggerBuilder {
	return &consoleLoggerBuilder{
		builder: NewGenericLoggerBuilder(),
	}
}

type consoleLoggerBuilder struct {
	builder *GenericLoggerBuilder
}

func (b *consoleLoggerBuilder) AddLevel(name string, display func() string) LoggerBuilder {
	b.builder.AddLevel(name, display)
	return b
}

func (b *consoleLoggerBuilder) AddPadding(padding Padding) LoggerBuilder {
	b.builder.AddPadding(padding)
	return b
}

func (b *consoleLoggerBuilder) AddColumn(column Column) LoggerBuilder {
	b.builder.AddColumn(column)
	return b
}

func (b *consoleLoggerBuilder) AddColumnByIndex(index uint, column Column) (LoggerBuilder, error) {
	err := b.builder.AddColumnByIndex(index, column)
	return b, err
}

func (b *consoleLoggerBuilder) Build() Logger {
	paddings := make([]Padding, 0)

	for key := range b.builder.paddings {
		paddings = append(paddings, key)
	}

	return &ConsoleLogger{
		levels:   b.builder.levels,
		paddings: paddings,
		columns:  b.builder.columns,
	}
}
