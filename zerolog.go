package golog

import (
  "github.com/rs/zerolog"
  "net"
  "time"
)

const (
  defaultZerologLevel = zerolog.InfoLevel
)

func NewZerologLogger(config Config) Interface {
  l := parseZerologLevel(config)

  writer, err := config.getOutputsWriter()
  if err != nil {
    panic(err)
  }

  logger := zerolog.New(writer).Level(l).With().Timestamp().Caller()

  for key, value := range config.Metadata {
    switch x := value.(type) {
    case string:
      logger = logger.Str(key, x)
    case int:
      logger = logger.Int(key, x)
    case bool:
      logger = logger.Bool(key, x)
    case net.IP:
      logger = logger.IPAddr(key, x)
    }
  }

  return &zerologLogger{
    config: config,
    logger: logger.Logger(),
  }
}

type zerologLogger struct {
  config Config
  logger zerolog.Logger
}

func (z *zerologLogger) NewWithFields(kv ...interface{}) Interface {
  context := z.logger.With()

  // Only want even num of key-value pairs
  if len(kv) % 2 == 0 {
    for i := 0; i < len(kv); i += 2 {
      if key, ok := kv[i].(string); ok {
        switch x := kv[i + 1].(type) {
        case string:
          context = context.Str(key, x)
        case int:
          context = context.Int(key, x)
        case bool:
          context = context.Bool(key, x)
        case error:
          context = context.AnErr(key, x)
        case time.Duration:
          context = context.Dur(key, x)
        }
      }
    }
  }

  return &zerologLogger{
    config: z.config,
    logger: context.Logger(),
  }
}

func (z *zerologLogger) Debug(msg string) {
  if e := z.logger.Debug(); e.Enabled() {
    e.Msg(msg)
  }
}

func (z *zerologLogger) Debugf(format string, args ...interface{}) {
  if e := z.logger.Debug(); e.Enabled() {
    e.Msgf(format, args...)
  }
}

func (z *zerologLogger) Info(msg string) {
  if e := z.logger.Info(); e.Enabled() {
    e.Msg(msg)
  }
}

func (z *zerologLogger) Infof(format string, args ...interface{}) {
  if e := z.logger.Info(); e.Enabled() {
    e.Msgf(format, args...)
  }
}

func (z *zerologLogger) Warn(msg string) {
  if e := z.logger.Warn(); e.Enabled() {
    e.Msg(msg)
  }
}

func (z *zerologLogger) Warnf(format string, args ...interface{}) {
  if e := z.logger.Warn(); e.Enabled() {
    e.Msgf(format, args...)
  }
}

func (z *zerologLogger) Error(msg string) {
  if e := z.logger.Error(); e.Enabled() {
    e.Msg(msg)
  }
}

func (z *zerologLogger) Errorf(format string, args ...interface{}) {
  if e := z.logger.Error(); e.Enabled() {
    e.Msgf(format, args...)
  }
}

func (z *zerologLogger) Fatal(msg string) {
  if e := z.logger.Fatal(); e.Enabled() {
    e.Msg(msg)
  }
}

func (z *zerologLogger) Fatalf(format string, args ...interface{}) {
  if e := z.logger.Fatal(); e.Enabled() {
    e.Msgf(format, args...)
  }
}

func parseZerologLevel(config Config) zerolog.Level {
  if l, err := zerolog.ParseLevel(config.Level); err != nil {
    return defaultZerologLevel
  } else {
    return l
  }
}
