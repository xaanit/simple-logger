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
	"github.com/logrusorgru/aurora/v3"
)

type LoggerBuilder interface {
	AddLevel(name string, display func() string) LoggerBuilder
	AddPadding(padding Padding) LoggerBuilder
	AddColumn(column Column) LoggerBuilder
	AddColumnByIndex(index int, column Column) LoggerBuilder
	Build() Logger
}

func defaultLevels(level string, display func() string, exclude []string, builder LoggerBuilder) {
	if findStrings(exclude, level) == -1 {
		builder.AddLevel(level, display)
	}
}

func defaultPaddings(padding Padding, exclude []Padding, builder LoggerBuilder) {
	if findPadding(exclude, padding) == -1 {
		builder.AddPadding(padding)
	}
}

func defaultColumns(column int, exclude []int, display func(context Context) string, builder LoggerBuilder) {
	if findInts(exclude, column) == -1 {
		builder.AddColumnByIndex(column, display)
	}
}

func SetDefaults(builder LoggerBuilder, excludeLevels []string, excludePaddings []Padding, excludeColumns []int) LoggerBuilder {
	defaultLevels("INFO", Info, excludeLevels, builder)
	defaultLevels("DEBUG", Debug, excludeLevels, builder)
	defaultLevels("ERROR", Error, excludeLevels, builder)
	defaultLevels("FATAL", Fatal, excludeLevels, builder)
	defaultLevels("WARNING", Warning, excludeLevels, builder)

	defaultPaddings(TimestampPadding, excludePaddings, builder)
	defaultPaddings(LevelPadding, excludePaddings, builder)

	defaultColumns(0, excludeColumns, func(context Context) string {
		layout := fmt.Sprintf("%v %v %v, %v | %v:%v:%v", Weekday, Month, Day, Year, Hour, Minute, Second)
		return aurora.BrightBlue(context.FormatTimestamp(layout)).String()
	}, builder)
	defaultColumns(1, excludeColumns, func(context Context) string { return context.FormatLevel() }, builder)
	defaultColumns(2, excludeColumns, func(context Context) string { return context.Message }, builder)

	return builder
}

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

func (b *GenericLoggerBuilder) AddLevel(name string, display func() string) {
	b.levels[name] = display
}

func (b *GenericLoggerBuilder) AddPadding(padding Padding) {
	b.paddings[padding] = nil
}

func (b *GenericLoggerBuilder) AddColumn(column Column) {
	b.columns = append(b.columns, column)
}

func (b *GenericLoggerBuilder) AddColumnByIndex(index int, column Column) {
	// https://stackoverflow.com/questions/46128016/insert-a-value-in-a-slice-at-a-given-index
	b.columns = append(b.columns, nil)
	copy(b.columns[index+1:], b.columns[index:])
	b.columns[index] = column
}
