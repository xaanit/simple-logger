package simple_logger

import (
	"fmt"
	"time"
)

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

func (c ConsoleLogger) GetLevels() map[string]func() string {
	return c.levels
}

func (c ConsoleLogger) GetPaddings() []Padding {
	return c.paddings
}

func (c ConsoleLogger) GetColumns() []Column {
	return c.columns
}

func (c ConsoleLogger) Log(level, message string) {
	if _, ok := c.GetLevels()[level]; !ok {
		return
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
}

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

func (b *consoleLoggerBuilder) AddColumnByIndex(index int, column Column) LoggerBuilder {
	b.builder.AddColumnByIndex(index, column)
	return b
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
