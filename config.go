package golog

import (
  "io"
  "os"
  "strings"

  _ "github.com/wgentry22/config"
)

const (
  ModuleName = `logger`
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

type Config struct {
  Level    string                 `mapstructure:"level"`
  Outputs  []string               `mapstructure:"outputs"`
  Metadata map[string]interface{} `mapstructure:"metadata"`
}

func (c *Config) getOutputsWriter() (io.Writer, error) {
  if len(c.Outputs) == 0 {
    return os.Stdout, nil
  }

  if len(c.Outputs) == 1 {
    return determineWriterFromOutput(c.Outputs[0])
  }

  outputErr := make(logOutputsMultiErr)
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
    if file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
      return nil, err
    } else {
      return file, nil
    }
  }
}
