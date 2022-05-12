<p align="center">
  <h3 align="center">logger</h3>
  <p align="center"><strong>logging utility for Golang</strong></p>

  <p align="center">
    <!-- Documentation -->
    <a href="http://godoc.org/github.com/fogfish/logger">
      <img src="https://godoc.org/github.com/fogfish/logger?status.svg" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/fogfish/logger/actions/">
      <img src="https://github.com/fogfish/logger/workflows/build/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/fogfish/logger">
      <img src="https://img.shields.io/github/last-commit/fogfish/logger.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/fogfish/logger?branch=main">
      <img src="https://coveralls.io/repos/github/fogfish/logger/badge.svg?branch=main" />
    </a>
    <!-- Go Card -->
    <a href="https://goreportcard.com/report/github.com/fogfish/logger">
      <img src="https://goreportcard.com/badge/github.com/fogfish/logger" />
    </a>
    <!-- Maintainability -->
    <a href="https://codeclimate.com/github/fogfish/logger/maintainability">
      <img src="https://api.codeclimate.com/v1/badges/df33ca9c2f9661803f78/maintainability" />
    </a>
  </p>
</p>

---

logger is a simple logging utility for Golang application. The library implements finer grained logging level, outputs log messages to standard streams. Easy to use for serverless application development (e.g. logging to AWS CloudWatch).


## Inspiration

The library outputs log messages in the well-defined format, using UTC date and times and giving ability to annotate the message with JSON context. The logger always output filename and line of the log statement to facilitate further analysis

```
2020/12/01 20:30:40 main.go:11: [level] some message
2020/12/01 20:30:40 main.go:11: [level] {"json": "context", "message": "some message"}
```

It inherits best practices of telecom application and defines 7 levels of fine grained logging:

1. `emergency`: system is unusable, panic execution of current routine or application, it is not possible to gracefully terminate it.
2. `critical`: system is failed, response actions must be taken immediately, the application is not able to execute correctly but still able to gracefully exit.
3. `error`: system is failed, unable to recover from error. The failure do not have global catastrophic impacts but local functionality is impaired, incorrect result is returned.
4. `warning`: system is failed, unable to recover, degraded functionality. The failure is ignored and application still capable to deliver incomplete but correct results.
5. `notice`: system is failed, error is recovered, no impact.
6. `info`: output informative status about system.
7. `debug`: output debug status about system.


## Getting started

The latest version of the library is available at `main` branch of this repository. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

Import library and start logging

```go
import "github.com/fogfish/logger"

logger.Notice("Some message")
logger.Error("Some formatted message %d", 5)
```

### Configuration

The default configuration applies debug level logging, messages are emitted to standard error (`os.Stderr`). Use `logger.Config` to change either logging level or destination. 

```go
logger.Config(logger.ERROR, os.Stderr)
```

Use environment variable `CONFIG_LOGGER_LEVEL` to change log level of the application at runtime

```bash
export CONFIG_LOGGER_LEVEL=warning
```

### Annotate events

Often, it is usable to annotate log messages with semi-structured data (e.g. JSON). 

```go
/*

Outputs

2020/12/01 20:30:40 main.go:11: [level] { "foo": "bar", "bar": 1, "message": "some message" }

*/
logger.With(logger.Note{
  "foo": "bar",
  "bar": 1
}).Warning("some message")
```

### AWS CloudWatch

The logger output events in the format compatible with AWS CloudWatch: each log message corresponds to single CloudWatch event. Therefore, it simplify logging in AWS Lambda functions. Use the logger together with CloudWatch Insight (e.g. utility [awslog](https://github.com/fogfish/awslog)) for the deep analysis. For example, search events with logs insight queries:

```
fields @timestamp, @message
| filter foo = "bar" and @message like /warning/
| sort @timestamp desc
| limit 20
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
go test
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/logger/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/logger.svg?style=for-the-badge)](LICENSE)
