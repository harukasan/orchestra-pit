package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

type formatterFunc func(e *Entry) []byte

func ColoredOutput(w io.Writer) Output {
	return &outputImpl{
		w:         w,
		formatter: formatColoredText,
	}
}

func TextOutput(w io.Writer) Output {
	return &outputImpl{
		w:         w,
		formatter: formatText,
	}
}

var StdoutOutput = ColoredOutput(colorable.NewColorableStdout())

func JSONOutput(w io.Writer) Output {
	return &outputImpl{
		w:         w,
		formatter: formatJSON,
	}
}

type outputImpl struct {
	w         io.Writer
	formatter formatterFunc
}

func (o *outputImpl) WriteEntry(e *Entry) {
	o.w.Write(o.formatter(e))
}

var (
	panicPrefix   = fmt.Sprintf("%-7s", "PANIC")
	fatalPrefix   = fmt.Sprintf("%-7s", "FATAL")
	errorPrefix   = fmt.Sprintf("%-7s", "ERROR")
	warningPrefix = fmt.Sprintf("%-7s", "WARN")
	noticePrefix  = fmt.Sprintf("%-7s", "NOTICE")
	infoPrefix    = fmt.Sprintf("%-7s", "INFO")
	debugPrefix   = fmt.Sprintf("%-7s", "DEBUG")
)

func formatText(e *Entry) []byte {
	buf := bytes.NewBuffer(nil)
	message := strings.TrimSpace(e.Message)

	buf.WriteString(e.Time.Format("2006-01-02T15:04:05-07:00"))
	buf.WriteRune(' ')

	switch e.Level {
	case PanicLevel:
		buf.WriteString(panicPrefix)
	case FatalLevel:
		buf.WriteString(fatalPrefix)
	case ErrorLevel:
		buf.WriteString(errorPrefix)
	case WarningLevel:
		buf.WriteString(warningPrefix)
	case NoticeLevel:
		buf.WriteString(noticePrefix)
	case InfoLevel:
		buf.WriteString(infoPrefix)
	case DebugLevel:
		buf.WriteString(debugPrefix)
	}

	buf.WriteString(message)
	buf.WriteRune('\n')

	return buf.Bytes()
}

var (
	coloredPanicPrefix   = color.New(color.FgRed, color.Bold).SprintfFunc()("%-7s", "PANIC")
	coloredFatalPrefix   = color.New(color.FgRed, color.Bold).SprintfFunc()("%-7s", "FATAL")
	coloredErrorPrefix   = color.New(color.FgRed).SprintfFunc()("%-7s", "ERROR")
	coloredWarningPrefix = color.New(color.FgYellow, color.Bold).SprintfFunc()("%-7s", "WARN")
	coloredNoticePrefix  = color.New(color.FgYellow).SprintfFunc()("%-7s", "NOTICE")
	coloredInfoPrefix    = color.New(color.FgCyan).SprintfFunc()("%-7s", "INFO")
	coloredDebugPrefix   = color.New(color.FgWhite).SprintfFunc()("%-7s", "DEBUG")
)

func formatColoredText(e *Entry) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(e.Time.Format("2006-01-02T15:04:05-07:00"))
	buf.WriteRune(' ')

	switch e.Level {
	case PanicLevel:
		buf.WriteString(coloredPanicPrefix)
	case FatalLevel:
		buf.WriteString(coloredFatalPrefix)
	case ErrorLevel:
		buf.WriteString(coloredErrorPrefix)
	case WarningLevel:
		buf.WriteString(coloredWarningPrefix)
	case NoticeLevel:
		buf.WriteString(coloredNoticePrefix)
	case InfoLevel:
		buf.WriteString(coloredInfoPrefix)
	case DebugLevel:
		buf.WriteString(coloredDebugPrefix)
	}

	message := strings.TrimSpace(e.Message)
	lines := strings.Split(message, "\n")
	buf.WriteString(lines[0])
	buf.WriteRune('\n')
	if len(lines) > 1 {
		for _, line := range lines[1:] {
			buf.WriteString("                                 ")
			buf.WriteString(line)
			buf.WriteRune('\n')
		}
	}

	return buf.Bytes()
}

func formatJSON(e *Entry) []byte {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return b
}
