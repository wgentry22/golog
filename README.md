## Golog
A common interface wrapping existing logging frameworks. Support currently exists for
* [Zerolog](https://github.com/rs/zerolog)
* Noop (useful for testing)


### Using `golog`
You get a `golog.Interface` instance by invoking `golog.NewLogger`, providing a `Config` instance

```go
package main

import "github.com/wgentry22/golog"

func main() {
  config := golog.Config{
    Level:    "debug",
    Outputs:  []string{"stdout", "path/to/my.log"},
    Metadata: map[string]interface{}{
      "service": "example",
    },
  }
  
  // Includes `service=example` on every log
  logger := golog.NewLogger(config)

  // Includes `service=example port=1234 env=qa` on every log
  subLogger := logger.NewWithFields("port", 1234, "env", "qa")
}
```

### Plays nicely with [`config`](https://github.com/wgentry22/config)
You can also provide a `config.Interface` instance to orchestrate creation of your `golog.Interface`
```go
package main

import (
  "github.com/wgentry22/config"
  "github.com/wgentry22/golog"
)

func main() {
  options := []config.Option{config.Name("myConfigFile"), config.Paths("path/to/my/config/")}
  conf := config.MustInit(options...)
  
  logger := golog.NewLoggerFromConfig(conf)
}
```
