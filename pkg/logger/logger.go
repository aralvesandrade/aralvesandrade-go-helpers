package logger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	color "github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

type ILogger interface {
	LogIt(severity, message string)
}

type Logger struct {
	logger *logrus.Logger
}

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	severity := colorLog(f.LevelDesc[entry.Level], f.LevelDesc[entry.Level])
	date := colorLog(f.LevelDesc[entry.Level], entry.Time.Format(f.TimestampFormat))
	message := colorLog(f.LevelDesc[entry.Level], entry.Message)
	formatter := fmt.Sprintf("%s %s %s\n", severity, date, message)
	return []byte(formatter), nil
}

func NewLogger() ILogger {
	g := &Logger{}
	g.logger = logrus.New()
	g.logger.SetLevel(getLogLevel(os.Getenv("LOGGING_LEVEL")))
	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = "02/01/2006 15:04:05"
	plainFormatter.LevelDesc = []string{"PANIC", "FAIL", "ERROR", "WARN", "INFO", "DEBUG"}
	g.logger.SetFormatter(plainFormatter)
	return g
}

func (l *Logger) LogIt(severity, message string) {
	details, _ := strconv.ParseBool(os.Getenv("LOGGING_DETAILS"))

	if details {
		pc, file, line, _ := runtime.Caller(1)
		fileName := fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
		funcName := runtime.FuncForPC(pc).Name()
		fn := funcName[strings.LastIndex(funcName, ".")+1:]
		message = fmt.Sprintf("(%s | %s) %s", fileName, fn, message)
	}

	switch severity {
	case "DEBUG":
		l.logger.Debug(message)
	case "INFO":
		l.logger.Info(message)
	case "WARN":
		l.logger.Warn(message)
	case "ERROR":
		l.logger.Error(message)
	default:
		l.logger.Info(message)
	}
}

func getLogLevel(env string) logrus.Level {
	switch env {
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func colorLog(severity, message string) string {
	message = strings.TrimSuffix(message, "\n")
	formattedmessage := fmt.Sprintf("[%v]", message)
	switch severity {
	case "DEBUG":
		return fmt.Sprint(color.Cyan(formattedmessage))
	case "INFO":
		return fmt.Sprint(color.Green(formattedmessage))
	case "WARN":
		return fmt.Sprint(color.Yellow(formattedmessage))
	case "ERROR":
		return fmt.Sprint(color.Red(formattedmessage))
	default:
		return fmt.Sprint(color.Green(formattedmessage))
	}
}
