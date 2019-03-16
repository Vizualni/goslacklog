# GoSlackLog

GoSlackLog is a simple library to log events to Slack. This is great for prototyping and should not be used on real production environments.

## Example

Image below shows simple log events from a test web server:

![GoSlackLog example](https://github.com/Vizualni/goslacklog/raw/master/example.png?sanitize=true)


## Example

```go
var slackLogger = goslacklog.NewBufferedSlackLogger(
	"https://hooks.slack.com/services/XXXXXXXXXXXXXX",
	1000, // buffer size of 1000 elements
	10, // send in 10 chunks
	2*time.Second, // every 2 seconds
)
var golog = log.New(slackLogger, "", log.LstdFlags)
golog.Printf("%s", "hello")
```
