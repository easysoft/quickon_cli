// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/easysoft/qcadmin/internal/pkg/util/log/survey"
)

// DiscardLogger just discards every log statement
type DiscardLogger struct {
	PanicOnExit bool
}

// Debug implements logger interface
func (d *DiscardLogger) Debug(args ...interface{}) {}

// Debugf implements logger interface
func (d *DiscardLogger) Debugf(format string, args ...interface{}) {}

// Info implements logger interface
func (d *DiscardLogger) Info(args ...interface{}) {}

// Infof implements logger interface
func (d *DiscardLogger) Infof(format string, args ...interface{}) {}

// Warn implements logger interface
func (d *DiscardLogger) Warn(args ...interface{}) {}

// Warnf implements logger interface
func (d *DiscardLogger) Warnf(format string, args ...interface{}) {}

// Error implements logger interface
func (d *DiscardLogger) Error(args ...interface{}) {}

// Errorf implements logger interface
func (d *DiscardLogger) Errorf(format string, args ...interface{}) {}

// Fatal implements logger interface
func (d *DiscardLogger) Fatal(args ...interface{}) {
	if d.PanicOnExit {
		d.Panic(args...)
	}

	os.Exit(1)
}

// Fatalf implements logger interface
func (d *DiscardLogger) Fatalf(format string, args ...interface{}) {
	if d.PanicOnExit {
		d.Panicf(format, args...)
	}

	os.Exit(1)
}

// Panic implements logger interface
func (d *DiscardLogger) Panic(args ...interface{}) {
	panic(fmt.Sprint(args...))
}

// Panicf implements logger interface
func (d *DiscardLogger) Panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

// Done implements logger interface
func (d *DiscardLogger) Done(args ...interface{}) {}

// Donef implements logger interface
func (d *DiscardLogger) Donef(format string, args ...interface{}) {}

// Print implements logger interface
func (d *DiscardLogger) Print(level logrus.Level, args ...interface{}) {}

// Printf implements logger interface
func (d *DiscardLogger) Printf(level logrus.Level, format string, args ...interface{}) {}

// StartWait implements logger interface
func (d *DiscardLogger) StartWait(message string) {}

// StopWait implements logger interface
func (d *DiscardLogger) StopWait() {}

// SetLevel implements logger interface
func (d *DiscardLogger) SetLevel(level logrus.Level) {}

// GetLevel implements logger interface
func (d *DiscardLogger) GetLevel() logrus.Level { return logrus.FatalLevel }

// Write implements logger interface
func (d *DiscardLogger) Write(message []byte) (int, error) {
	return len(message), nil
}

// WriteString implements logger interface
func (d *DiscardLogger) WriteString(message string) {}

// Question asks a new question
func (d *DiscardLogger) Question(params *survey.QuestionOptions) (string, error) {
	return "", SurveyError{}
}

// SurveyError is used to identify errors where questions were asked in the discard logger
type SurveyError struct{}

// Error implements error interface
func (s SurveyError) Error() string {
	return "Asking questions is not possible in silenced mode"
}
