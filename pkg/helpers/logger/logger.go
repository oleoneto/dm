package logger

import (
	"fmt"
	"io"
	"strings"

	"github.com/drewstinnett/go-output-format/formatter"
)

type Logger struct {
	config         formatter.Config
	cachedMessages []Formattable
}

type Formattable interface {
	Description() string
}

type ApplicationMessage struct {
	Message string `json:"message" yaml:"message"`
}

type ApplicationMessages []ApplicationMessage

type ApplicationError struct {
	Error string `json:"error" yaml:"error"`
}

func (l *Logger) CacheMessage(message Formattable /*, ioWritter io.Writer*/) {
	l.cachedMessages = append(l.cachedMessages, message)
}

func (l Logger) WithFormattedOutput(data Formattable, ioWriter io.Writer) {
	output, err := formatter.OutputData(data, &l.config)

	if err != nil {
		fmt.Fprintln(ioWriter, err)
		return
	}

	switch l.config.Format {
	case "plain":
		fmt.Fprintln(ioWriter, data.Description())
	default:
		fmt.Fprintln(ioWriter, string(output))
	}
}

func (l *Logger) ReleaseCachedMessages(ioWriter io.Writer) {
	switch l.config.Format {
	case "plain":
		messages := []string{}

		for _, m := range l.cachedMessages {
			messages = append(messages, m.Description())
		}

		fmt.Fprintln(ioWriter, strings.Join(messages, "\n"))
	default:
		output, err := formatter.OutputData(l.cachedMessages, &l.config)

		if err != nil {
			fmt.Fprintln(ioWriter, err)
			return
		}

		fmt.Fprintln(ioWriter, string(output))
	}
}

func (t ApplicationError) Description() string {
	return t.Error
}

func (t ApplicationMessage) Description() string {
	return t.Message
}

// MARK: Default Loggers

func Default() Logger {
	return Logger{
		config: formatter.Config{
			Template: "",
			Format:   "plain",
		},
	}
}

func Custom(format, template string) Logger {
	return Logger{
		config: formatter.Config{
			Template: template,
			Format:   format,
		},
	}
}
