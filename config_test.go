package golog

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/spf13/afero"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"
  "io"
  "os"
  "reflect"
  "testing"
)

const (
  forbiddenDir = `/forbidden`
)

type ConfigTestSuite struct {
  suite.Suite
}

func (c *ConfigTestSuite) SetupSuite() {
  testFileSystem := afero.NewMemMapFs()
  err := testFileSystem.Mkdir(forbiddenDir, 0000)
  require.Nil(c.T(), err)
}

func (c *ConfigTestSuite) TestCreateNoopLogger() {
  config := Config{
    Kind:     "noop",
    Level:    "debug",
    Outputs:  []string{"stdout"},
    Metadata: make(map[string]interface{}),
  }

  logger := NewLogger(config)
  require.NotNil(c.T(), logger)
  require.Equal(c.T(), reflect.TypeOf(logger), reflect.TypeOf(&noopLogger{}))
}

func (c *ConfigTestSuite) TestCreateZerologLogger() {
  config := Config{
    Level:    "debug",
    Outputs:  []string{"stdout"},
    Metadata: make(map[string]interface{}),
  }

  logger := NewLogger(config)
  require.NotNil(c.T(), logger)
  require.Equal(c.T(), reflect.TypeOf(logger), reflect.TypeOf(&zerologLogger{}))
}

func (c *ConfigTestSuite) TestGetWriters() {
  tempDir := c.T().TempDir()
  tempFileName := fmt.Sprintf("%s/%s.log", tempDir, uuid.New())

  tempFile, err := os.OpenFile(tempFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
  require.Nil(c.T(), err)
  require.NotNil(c.T(), tempFile)

  config := Config{
    Level:    "debug",
    Outputs:  []string{"stdout", "stderr", tempFileName},
    Metadata: make(map[string]interface{}),
  }

  writer, err := config.getOutputsWriter()
  require.Nil(c.T(), err)
  require.NotNil(c.T(), writer)
  require.Equal(c.T(), reflect.TypeOf(writer), reflect.TypeOf(io.MultiWriter(os.Stdout, os.Stderr, tempFile)))
}

func (c *ConfigTestSuite) TestGetWriters_ReturnsStdout_WhenNoOutputsProvided() {
  config := Config{
    Level:    "debug",
    Metadata: make(map[string]interface{}),
  }

  writer, err := config.getOutputsWriter()
  require.Nil(c.T(), err)
  require.NotNil(c.T(), writer)
  require.Equal(c.T(), writer, os.Stdout)
}

func (c *ConfigTestSuite) TestGetWriters_ReturnsErr_WhenErrorOpeningFile() {
  forbiddenFile := fmt.Sprintf("%s/%s.log", forbiddenDir, uuid.New())

  config := Config{
    Outputs: []string{forbiddenFile},
  }

  writer, err := config.getOutputsWriter()
  require.Nil(c.T(), writer)
  require.NotNil(c.T(), err)
  require.Equal(c.T(), reflect.TypeOf(err), reflect.TypeOf(&logOutputsMultiErr{}))
}

func (c *ConfigTestSuite) TestGetWriters_ReturnsErr_WhenErrorOpeningMultipleFiles() {
  forbiddenFile := fmt.Sprintf("%s/%s.log", forbiddenDir, uuid.New())

  config := Config{
    Outputs: []string{"stdout", forbiddenFile},
  }

  writer, err := config.getOutputsWriter()
  require.Nil(c.T(), writer)
  require.NotNil(c.T(), err)
  require.Equal(c.T(), reflect.TypeOf(err), reflect.TypeOf(&logOutputsMultiErr{}))
}

func TestConfigTestSuite(t *testing.T) {
  suite.Run(t, new(ConfigTestSuite))
}
