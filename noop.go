// Copyright 2022 wgentry22. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golog

func newNoopLogger(config Config) Interface {
  return &noopLogger{}
}

type noopLogger struct {}

func (n *noopLogger) NewWithFields(kv ...interface{}) Interface {
  return n
}

func (n *noopLogger) Debug(msg string) {}

func (n *noopLogger) Debugf(format string, args ...interface{}) {}

func (n *noopLogger) Info(msg string) {}

func (n *noopLogger) Infof(format string, args ...interface{}) {}

func (n *noopLogger) Warn(msg string) {}

func (n *noopLogger) Warnf(format string, args ...interface{}) {}

func (n *noopLogger) Error(msg string) {}

func (n *noopLogger) Errorf(format string, args ...interface{}) {}

func (n *noopLogger) Fatal(msg string) {}

func (n *noopLogger) Fatalf(format string, args ...interface{}) {}



