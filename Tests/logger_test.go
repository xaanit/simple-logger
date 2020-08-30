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
package Tests

import (
	"fmt"
	log "simple-logger"
	"testing"
	"time"
)

type testLogger struct {
	levels   map[string]func() string
	paddings []log.Padding
	columns  []log.Column
	str      string
}

func (t *testLogger) GetLevels() map[string]func() string {
	return t.levels
}

func (t testLogger) GetPaddings() []log.Padding {
	return t.paddings
}

func (t testLogger) GetColumns() []log.Column {
	return t.columns
}

func (c *testLogger) createContext(level, message string) log.Context {
	return log.Context{
		Message: message,
		Time:    time.Now(),
		Level:   level,
		Logger:  c,
	}
}

func (c *testLogger) Log(level, message string) (int, error) {
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

	c.str = format

	return log.Success, nil
}

func TestLogger(t *testing.T) {
	builder := log.NewGenericLoggerBuilder()
	log.SetDefaults(builder, nil, nil, []uint{0})
	paddings := make([]log.Padding, 0)
	for padding := range builder.Paddings {
		paddings = append(paddings, padding)
	}
	logger := &testLogger{
		levels:   builder.Levels,
		paddings: paddings,
		columns:  builder.Columns,
		str:      "",
	}

	_, _ = logger.Log("INFO", "Hello, world!")
	expected := fmt.Sprintf("%v    | Hello, world!", log.Info())
	if logger.str != expected {
		t.Fatalf("Logger.str was [%v] not [%v]", logger.str, expected)
	}
}
