# Application Formatter
![Test Status](https://github.com/SivWatt/formatter/actions/workflows/go.yaml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/SivWatt/formatter)](https://goreportcard.com/report/github.com/SivWatt/formatter)  
This is a customized formatter which implements `Fommatter` of [logrus](https://github.com/sirupsen/logrus).  
This formatter mainly focuses on desktop application logging, since we might need __process ID__, __local time stamp__, __function name__, and __lines__.  
And one of the goals of this formatter is making logs __human-readable__.  
```log
2021-04-25T20:37:27+08:00 [TRAC] [2070] some trace message
2021-04-25T20:37:27+08:00 [DEBU] [2070] some debug message
2021-04-25T20:37:27+08:00 [INFO] [2070] some info message
2021-04-25T20:37:27+08:00 [WARN] [2070] some warning message
2021-04-25T20:37:27+08:00 [ERRO] [2070] some error message
2021-04-25T20:37:27+08:00 [FATA] [2070] some fatal message
```

## Configuration
```golang
type AppFormatter struct {
	DisableTimestamp bool
	DisablePID       bool
}
```

## Usage
```golang
import (
	"github.com/SivWatt/formatter"
	"github.com/sirupsen/logrus"
)

log := logrus.New()
log.SetFormatter(&formatter.AppFormatter{})

log.Info("some info message")
// Output: 2021-04-25T20:30:03+08:00 [INFO] [1228] some info message

log.WithField("key", "value").Debug("some debug message")
// Output: 2021-04-25T20:30:03+08:00 [DEBU] [1228] some debug message [key:value]
```

## Development
Feel free to let me know if there is anything needed to be improved. _(via issue)_