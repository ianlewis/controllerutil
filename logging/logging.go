// Copyright 2017 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logging provides standard log.Logger wrappers around the github.com/golang/glog package.
// These wrappers can be set with log message prefixes.
package logging

import (
	"log"

	"github.com/golang/glog"
)

type infoWriter struct {
	level glog.Level
}

// Write implements the io.Writer interface.
func (w infoWriter) Write(data []byte) (n int, err error) {
	glog.V(w.level).Info(string(data))
	return len(data), nil
}

type errorWriter struct{}

// Write implements the io.Writer interface.
func (w errorWriter) Write(data []byte) (n int, err error) {
	glog.Error(string(data))
	return len(data), nil
}

// Logger is a simple wrapper around glog that provides a way to easily create loggers at a particular log level.
type Logger struct {
	prefix string
}

// V returns a standard log.Logger at the particular log level. Logs will only be printed if glog was set up at a level equal to or higher than the level provided to V.
func (l *Logger) V(level glog.Level) *log.Logger {
	return log.New(infoWriter{level}, l.prefix, 0)
}

// NewInfoLogger returns a Logger object that provides a V method to log info messages at a particular log level.
func NewInfoLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
	}
}

// NewErrorLogger returns a standard log.Logger that can be used to write error log messages.
func NewErrorLogger(prefix string) *log.Logger {
	return log.New(errorWriter{}, prefix, 0)
}
