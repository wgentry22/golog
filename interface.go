package golog

type Constructor func(Config) Interface

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
