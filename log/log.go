/*
 * @Descripttion:
 * @version:
 * @Author: 江小凡
 * @Date: 2022-01-26 22:28:51
 * @LastEditTime: 2022-01-26 22:46:23
 */
package log

import (
	"io"
	"mongo/conf"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

type Logger struct {
	Logger *logrus.Logger
}

// 创建一个日志记录器
func NewLogger() *Logger {
	l := logrus.New()
	return &Logger{l}
}

// 设置日志级别
func (l *Logger) SetLevel(level string) {
	switch level {
	case "debug":
		l.Logger.SetLevel(logrus.DebugLevel)
	case "info":
		l.Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		l.Logger.SetLevel(logrus.WarnLevel)
	case "error":
		l.Logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		l.Logger.SetLevel(logrus.FatalLevel)
	case "panic":
		l.Logger.SetLevel(logrus.PanicLevel)
	}
}

// 设置日志格式
func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.Formatter = &logrus.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableSorting:   false,
		QuoteEmptyFields: true,
		FieldMap:         nil,
		CallerPrettyfier: nil,
	}
}

// 设置日志输出
func (l *Logger) SetOutput(out io.Writer) {
	// 支持多个输出 out = io.MultiWriter(wirte1,write2)
	l.Logger.SetOutput(out)
}

// 设置日志输出文件名
func (l *Logger) SetFile(file string) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l.Logger.SetOutput(f)
}

// 添加日志输出字段
func (l *Logger) AddField(key string, value interface{}) {
	l.Logger.WithFields(logrus.Fields{key: value})
}

// 重定向日志输出
func (l *Logger) Redirect(w io.Writer) {
	l.Logger.Out = w
}

func init() {
	logger := NewLogger()
	logger.SetLevel(conf.LOG_LEVEL)
	logger.SetFormatter(&logrus.JSONFormatter{})
	Log = logger.Logger
	Log.Formatter = &logrus.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableSorting:   false,
		QuoteEmptyFields: true,
		FieldMap:         nil,
		CallerPrettyfier: nil,
	}

}
