// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log/survey"
)

var logs = map[string]Logger{}
var logsMutext sync.Mutex

var overrideOnce sync.Once

type fileLogger struct {
	logger *logrus.Logger
}

// GetFileLogger returns a logger instance for the specified filename
func GetFileLogger(filename string) Logger {
	filedir := common.GetDefaultLogDir()
	filename = strings.TrimSpace(filename)

	logsMutext.Lock()
	defer logsMutext.Unlock()

	log := logs[filename]
	if log == nil {
		newLogger := &fileLogger{
			logger: logrus.New(),
		}
		newLogger.logger.Formatter = &logrus.JSONFormatter{}
		newLogger.logger.SetOutput(&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s", filedir, filename),
			MaxAge:     12,
			MaxBackups: 4,
			MaxSize:    10 * 1024 * 1024,
		})

		newLogger.SetLevel(GetInstance().GetLevel())
		logs[filename] = newLogger
	}

	return logs[filename]
}

// OverrideRuntimeErrorHandler overrides the standard runtime error handler that logs to stdout
// with a file logger that logs all runtime.HandleErrors to errors.log
func OverrideRuntimeErrorHandler(discard bool) {
	overrideOnce.Do(func() {
		if discard {
			if len(runtime.ErrorHandlers) > 0 {
				runtime.ErrorHandlers[0] = func(_ context.Context, err error, msg string, keysAndValues ...interface{}) {}
			} else {
				runtime.ErrorHandlers = []runtime.ErrorHandler{
					func(_ context.Context, err error, msg string, keysAndValues ...interface{}) {},
				}
			}
		} else {
			errorLog := GetFileLogger("errors")
			if len(runtime.ErrorHandlers) > 0 {
				runtime.ErrorHandlers[0] = func(_ context.Context, err error, msg string, keysAndValues ...interface{}) {
					errorLog.Errorf("Runtime error occurred: %s", err)
				}
			} else {
				runtime.ErrorHandlers = []runtime.ErrorHandler{
					func(_ context.Context, err error, msg string, keysAndValues ...interface{}) {
						errorLog.Errorf("Runtime error occurred: %s", err)
					},
				}
			}
		}
	})
}

func (f *fileLogger) Debug(args ...interface{}) {
	f.logger.Debug(args...)
}

func (f *fileLogger) Debugf(format string, args ...interface{}) {
	f.logger.Debugf(format, args...)
}

func (f *fileLogger) Info(args ...interface{}) {
	f.logger.Info(args...)
}

func (f *fileLogger) Infof(format string, args ...interface{}) {
	f.logger.Infof(format, args...)
}

func (f *fileLogger) Warn(args ...interface{}) {
	f.logger.Warn(args...)
}

func (f *fileLogger) Warnf(format string, args ...interface{}) {
	f.logger.Warnf(format, args...)
}

func (f *fileLogger) Error(args ...interface{}) {
	f.logger.Error(args...)
}

func (f *fileLogger) Errorf(format string, args ...interface{}) {
	f.logger.Errorf(format, args...)
}

func (f *fileLogger) Fatal(args ...interface{}) {
	f.logger.Fatal(args...)
}

func (f *fileLogger) Fatalf(format string, args ...interface{}) {
	f.logger.Fatalf(format, args...)
}

func (f *fileLogger) Panic(args ...interface{}) {
	f.logger.Panic(args...)
}

func (f *fileLogger) Panicf(format string, args ...interface{}) {
	f.logger.Panicf(format, args...)
}

func (f *fileLogger) Done(args ...interface{}) {
	f.logger.Info(args...)
}

func (f *fileLogger) Donef(format string, args ...interface{}) {
	f.logger.Infof(format, args...)
}

func (f *fileLogger) Print(level logrus.Level, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		f.Info(args...)
	case logrus.DebugLevel:
		f.Debug(args...)
	case logrus.WarnLevel:
		f.Warn(args...)
	case logrus.ErrorLevel:
		f.Error(args...)
	case logrus.PanicLevel:
		f.Panic(args...)
	case logrus.FatalLevel:
		f.Fatal(args...)
	}
}

func (f *fileLogger) Printf(level logrus.Level, format string, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		f.Infof(format, args...)
	case logrus.DebugLevel:
		f.Debugf(format, args...)
	case logrus.WarnLevel:
		f.Warnf(format, args...)
	case logrus.ErrorLevel:
		f.Errorf(format, args...)
	case logrus.PanicLevel:
		f.Panicf(format, args...)
	case logrus.FatalLevel:
		f.Fatalf(format, args...)
	}
}

func (f *fileLogger) StartWait(message string) {
	// Noop operation
}

func (f *fileLogger) StopWait() {
	// Noop operation
}

func (f *fileLogger) SetLevel(level logrus.Level) {
	f.logger.SetLevel(level)
}

func (f *fileLogger) GetLevel() logrus.Level {
	return f.logger.GetLevel()
}

func (f *fileLogger) Write(message []byte) (int, error) {
	return f.logger.Out.Write(message)
}

func (f *fileLogger) WriteString(message string) {
	f.logger.Info(strings.TrimSuffix(message, "\n"))
}

func (f *fileLogger) Question(params *survey.QuestionOptions) (string, error) {
	return "", errors.New("questions in file logger not supported")
}
