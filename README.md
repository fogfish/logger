<p align="center">
  <h3 align="center">logger</h3>
  <p align="center"><strong>logging config for Golang's slog</strong></p>

  <p align="center">
    <!-- Documentation -->
    <a href="http://godoc.org/github.com/fogfish/logger">
      <img src="https://godoc.org/github.com/fogfish/logger?status.svg" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/fogfish/logger/actions/">
      <img src="https://github.com/fogfish/logger/workflows/test/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/fogfish/logger">
      <img src="https://img.shields.io/github/last-commit/fogfish/logger.svg" />
    </a>
  </p>
</p>

---

`logger` is an opinionated configuration for Golang [slog](https://pkg.go.dev/log/slog) developed for easy to use within serverless applications, such as logging to AWS CloudWatch.


## Inspiration

The library enables 0-configuration for [slog](https://pkg.go.dev/log/slog) enabling two logging formats:
* Log messages as JSON object compatible with CloudWatch for production.
* Log messages as colored text and JSON for development

```
2023-10-26 19:12:44.709 +0300 EEST:
{
    "level": "INFO",
    "source": {
        "function": "main.main",
        "file": "gthb.fgfs.lggr.exmp/main.go",
        "line": 143
    },
    "msg": "informative status about system.",
    "key": "val",
    ...
}
```

```
[14:06:54.030] INF informative status about system. {
  "key": "val",
  "source": {
    "file": "gthb.fgfs.lggr.exmp/main.go",
    "function": "main.main",
    "line": 26
  }
}
```

Additionally, it provides several enhancements:
* 7-level logging semantics for more precise error handling.
* File name shortening for cleaner and more readable logs.
* Module-based log level configuration for flexible logging control.
* Configuration via environment variables for easy customization.
* Built-in metrics for improved observability.


## Getting started

The latest version of the configuration is available at `main` branch of this repository. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

Import configuration and start logging using `slog` api. The default config is optimized for logging within Serverless application.

```bash
go get -u github.com/fogfish/logger/v3
```

- [Inspiration](#inspiration)
- [Getting started](#getting-started)
  - [Quick Start](#quick-start)
  - [Extended Logging Levels](#extended-logging-levels)
  - [Configuration](#configuration)
  - [Module-Based Log Level Configuration](#module-based-log-level-configuration)
  - [AWS CloudWatch](#aws-cloudwatch)
  - [Observability metrics](#observability-metrics)
- [How To Contribute](#how-to-contribute)
  - [commit message](#commit-message)
  - [bugs](#bugs)
- [License](#license)


### Quick Start

```go
package main

import (
  "context"
  "log/slog"

  log "github.com/fogfish/logger/v3"
)

func main() {
  slog.SetDefault(log.New())

  slog.Info("informative status about system.")
  slog.Warn("system is failed, unable to recover, degraded functionality.")
	slog.Error("system is failed, unable to recover from error.")
}
```

### Extended Logging Levels 

The `logger` library supplies configuration follows best practices from telecom applications, enhancing the standard `Debug`, `Info`, `Warn`, and `Error` levels with three additional levels. This provides fine-grained control over logging, resulting in seven distinct levels for correct error handling: 
1. `EMERGENCY` (`EMR`) – The system is unusable. A panic occurs, and it is impossible to gracefully terminate the application.
2. `CRITICAL` (`CRT`) – The system has failed and requires immediate action. The application cannot function correctly but can still exit gracefully.
3. `ERROR` (`ERR`) – A failure has occurred, and recovery is not possible. While the issue does not have catastrophic global effects, local functionality is impaired, leading to incorrect results.
4. `WARN` (`WRN`) – A failure has occurred, and recovery is not possible. However, the system continues operating in a degraded state, delivering incomplete but correct results.
5. `NOTICE` (`NTC`) – A failure occurred but was successfully recovered, with no lasting impact on the system.
6. `INFO` (`INF`) – Provides informational updates on the system’s status.
7. `DEBUG` (`DEB`) – Outputs detailed debugging information for troubleshooting.

This structured logging approach ensures clear categorization of system states, making it easier to detect, react to, and diagnose issues in complex applications. 

The faster way apply these levels is raw `slog.Log` function and standartized constants:  

```go
import (
  log "github.com/fogfish/logger/v3"
)

slog.Log(context.Background(), log.DEBUG, "...")
slog.Log(context.Background(), log.INFO, "...")
slog.Log(context.Background(), log.NOTICE, "...")
slog.Log(context.Background(), log.WARN, "...")
slog.Log(context.Background(), log.ERROR, "...")
slog.Log(context.Background(), log.CRITICAL, "...")
slog.Log(context.Background(), log.EMERGENCY, "...")
```

Alternatively, module `xlog` provides variants of these functions:

```go
import (
  "github.com/fogfish/logger/x/xlog"
)

xlog.Notice("...")
xlog.Warn("...", err)
xlog.Error("...", err)
xlog.Critical("...", err)
xlog.Emergency("...", err)
```

### Configuration

The typical configuration is following:

```go
import (
  "log/slog"
  log "github.com/fogfish/logger/v3"
)

slog.SetDefault(log.New())
```

The default configuration works out-of-the-box, automatically adapting to the runtime environment. Adjust it Using functional option pattern, see all configuration options and presets [here](./options.go).

The default log level is `INFO` and log messages are emitted to standard error (`os.Stderr`). Use environment variable `CONFIG_LOG_LEVEL` to change log level of the application at runtime:

```bash
export CONFIG_LOGGER_LEVEL=WARN
```

Note: the environemnt configuration is case sensitive, all caps is required. 


### Module-Based Log Level Configuration

The logger allows you to define log levels for different modules with flexible granularity. Log levels can be set explicitly using configuration options or environment variables.

The logger uses prefix matching to determine the appropriate log level based on the source code path:

**Per File**: A log level defined for a specific file (e.g., `github.com/fogfish/logger/logger.go`) applies only to that file.

**Per Module**: A log level set for a module (e.g., `github.com/fogfish/logger`) applies to all files within that module.

**Per Namespace**: A log level defined at a higher level (e.g., `github.com/fogfish`) applies to all modules under that namespace.

You either do explicit configuration using the config option

```go
import (
  "log/slog"

  log "github.com/fogfish/logger/v3"
)

slog.SetDefault(
  log.New(
    log.WithLogLevelForMod(map[string]slog.Level{
      "github.com/fogfish/logger":  log.INFO,
      "github.com/you/application": log.DEBUG,
    }),
  ),
)
```

Or, using environment variable `CONFIG_LOG_LEVEL_{LEVEL_NAME}`

```bash
export CONFIG_LOG_LEVEL_DEBUG=github.com/you/application:github.com/
export CONFIG_LOG_LEVEL_INFO=github.com/fogfish/logger
```


### AWS CloudWatch

The logger output events in the format compatible with AWS CloudWatch: each log message corresponds to single CloudWatch event. Therefore, it simplify logging in AWS Lambda functions. Use the logger together with CloudWatch Insight (e.g. utility [awslog](https://github.com/fogfish/awslog)) for the deep analysis. For example, search events with logs insight queries:

```
fields @timestamp, @message
| filter level = "INFO" and foo = "bar"
| sort @timestamp desc
| limit 20
```

### Observability metrics

Logging **duration** of the function

```go
func do() {
  defer slog.Info("done something", slog.Any("duration", xlog.SinceNow()))
  // ...
}
```

Logging **execution rate** of the code block.

```go
func do() {
  ops := xlog.PerSecondNow()
  defer slog.Info("done something", slog.Any("op/sec", ops))
  // ops.Acc++
}
```

Logging **demand** of the code block

```go
func do() {
  ops := xlog.MillisecondOpNow()
  defer slog.Info("done something", slog.Any("op/sec", ops))
  // ops.Acc++
}
```

## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.16 or later.

**build** and **test** library.

```bash
git clone https://github.com/fogfish/logger
cd logger
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/logger/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/logger.svg?style=for-the-badge)](LICENSE)
