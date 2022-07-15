// Copyright 2022 wgentry22. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golog

import (
  "github.com/spf13/afero"
  "io"
  "os"
  "strings"

  "github.com/wgentry22/config"
)

const (
  // ModuleName uniquely identifies the name of this package
  // When used with config.Config, ModuleName should be the key provided to config.Config#Get
  ModuleName = `logger`
)

var (
  appFileSystem = afero.NewOsFs()
)

type logOutputsMultiErr map[string]error

func (l *logOutputsMultiErr) Error() string {
  var builder strings.Builder
  builder.WriteString("Encountered errors getting output: ")
  for output, err := range *l {
    builder.WriteString("\n\t-")
    builder.WriteString(output)
    builder.WriteString(" => ")
    builder.WriteString(err.Error())
  }

  return builder.String()
}

// Config provides the possible options to customize the Interface instance
type Config struct {
  // Kind refers to the type of supported logger: i.e., `noop`, `zerolog`
  Kind     string                 `mapstructure:"kind"`
  // Level is desired the logging level
  Level    string                 `mapstructure:"level"`
  // Outputs are filesystem paths to which logs should be written to. Also accepts `stdout` and `stderr`.
  Outputs  []string               `mapstructure:"outputs"`
  // Metadata are key-value pairs that are included on every log
  Metadata map[string]interface{} `mapstructure:"metadata"`
}

// Interface provides a common logger interface agnostic of the underlying implementation
type Interface interface {
  NewWithFields(kv ...interface{}) Interface
  Debug(msg string)
  Debugf(format string, args ...interface{})
  Info(msg string)
  Infof(format string, args ...interface{})
  Warn(msg string)
  Warnf(format string, args ...interface{})
  Error(msg string)
  Errorf(format string, args ...interface{})
  Fatal(msg string)
  Fatalf(format string, args ...interface{})
}

// NewLogger create an Interface instance based on the provided Config
func NewLogger(config Config) Interface {
  switch config.Kind {
  case "noop":
    return newNoopLogger(config)
  default:
    return newZerologLogger(config)
  }
}

func NewLoggerFromConfig(config config.Interface) Interface {
  var c Config
  if err := config.Get(ModuleName, &c); err != nil {
    panic(err)
  }

  return NewLogger(c)
}

func (c *Config) getOutputsWriter() (io.Writer, error) {
  if len(c.Outputs) == 0 {
    return os.Stdout, nil
  }

  outputErr := make(logOutputsMultiErr)

  if len(c.Outputs) == 1 {
    writer, err := determineWriterFromOutput(c.Outputs[0])
    if err != nil {
      outputErr[c.Outputs[0]] = err
      return nil, &outputErr
    }

    return writer, nil
  }

  writers := make([]io.Writer, 0, len(c.Outputs))

  for _, output := range c.Outputs {
    writer, err := determineWriterFromOutput(output)
    if err != nil {
      outputErr[output] = err
    } else {
      writers = append(writers, writer)
    }
  }

  if len(outputErr) > 0 {
    return nil, &outputErr
  }

  return io.MultiWriter(writers...), nil
}

func determineWriterFromOutput(output string) (io.Writer, error) {
  switch output {
  case "stdout":
    return os.Stdout, nil
  case "stderr":
    return os.Stderr, nil
  default:
    if file, err := appFileSystem.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
      return nil, err
    } else {
      return file, nil
    }
  }
}
