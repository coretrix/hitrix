package errorlogger

import (
	"bytes"
	//nolint //G501: Blocklisted import crypto/md5: weak cryptographic primitive
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httputil"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"
	slackgo "github.com/slack-go/slack"

	"github.com/coretrix/hitrix/service/component/app"
	requestlogger "github.com/coretrix/hitrix/service/component/request_logger"
	"github.com/coretrix/hitrix/service/component/sentry"
	"github.com/coretrix/hitrix/service/component/slack"
)

const (
	GroupError              = "error"
	GroupWarning            = "warning"
	GroupMissingTranslation = "missing-translation"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type ErrorLogger interface {
	LogError(dataFromRecover interface{})
	LogErrorWithRequest(c *gin.Context, errData interface{})
	LogPanicWithRequest(c *gin.Context, errData interface{})
	LogWarning(data interface{})
	LogWarningWithRequest(c *gin.Context, data interface{})
	LogMissingTranslation(data interface{})
}

type ErrorMessage struct {
	File    string
	Line    int
	AppName string
	Request []byte
	Message string
	Stack   []byte
}

type RedisErrorLogger struct {
	redisStorage   *beeorm.RedisCache
	sentryService  sentry.ISentry
	slackService   slack.Slack
	appService     *app.App
	requestBodyKey interface{}
}

func NewRedisErrorLogger(
	appService *app.App,
	ormService *beeorm.Engine,
	slackService slack.Slack,
	sentryService sentry.ISentry,
	requestBodyKey interface{},
) ErrorLogger {
	redisStorage := ormService.GetRedis()
	if appService.RedisPools != nil && appService.RedisPools.Persistent != "" {
		redisStorage = ormService.GetRedis(appService.RedisPools.Persistent)
	}

	return &RedisErrorLogger{
		redisStorage:   redisStorage,
		slackService:   slackService,
		appService:     appService,
		sentryService:  sentryService,
		requestBodyKey: requestBodyKey,
	}
}

func (e *RedisErrorLogger) LogError(errData interface{}) {
	e.log(errData, nil, GroupError, debug.Stack())
}

func (e *RedisErrorLogger) LogErrorWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, c, GroupError, debug.Stack())
}

func (e *RedisErrorLogger) LogPanicWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, c, GroupError, debug.Stack())
}

func (e *RedisErrorLogger) LogWarning(data interface{}) {
	e.log(data, nil, GroupWarning, nil)
}

func (e *RedisErrorLogger) LogWarningWithRequest(c *gin.Context, data interface{}) {
	e.log(data, c, GroupWarning, nil)
}

func (e *RedisErrorLogger) LogMissingTranslation(data interface{}) {
	e.log(data, nil, GroupMissingTranslation, nil)
}

// stackFileLineRe matches the file:line line in debug.Stack() output (tab-indented path).
var stackFileLineRe = regexp.MustCompile(`^\s+(.+\.go):(\d+)\s+`)

// parseStackFileLine parses debug.Stack() output and returns file and line for the error location.
// Skips frames inside this package (error_logger). When the caller is a defer after recover(),
// the first external frame is the defer; the second is the panic site. So we use frame[1] when
// we have 2+ external frames, otherwise frame[0].
func parseStackFileLine(stack []byte) (file string, line int) {
	var external []struct {
		file string
		line int
	}
	for _, lineBytes := range bytes.Split(stack, []byte{'\n'}) {
		if m := stackFileLineRe.FindSubmatch(lineBytes); len(m) == 3 {
			f := string(m[1])
			if strings.Contains(f, "error_logger") {
				continue
			}
			var l int
			_, _ = fmt.Sscanf(string(m[2]), "%d", &l)
			external = append(external, struct {
				file string
				line int
			}{f, l})
		}
	}
	if len(external) == 0 {
		return "", 0
	}
	idx := 0
	if len(external) >= 2 {
		idx = 1 // panic path: first external is defer, second is panic site
	}
	return external[idx].file, external[idx].line
}

func (e *RedisErrorLogger) log(errData interface{}, c *gin.Context, group string, capturedStack []byte) {
	var msg string

	switch v := errData.(type) {
	case error:
		msg = v.Error()
	case string:
		msg = v
	default:
		msg = fmt.Sprint(v)
	}

	var stackTrace []byte
	var file string
	var line int

	if len(capturedStack) > 0 {
		stackTrace = capturedStack
		file, line = parseStackFileLine(capturedStack)
	} else {
		stackTrace = stack(0)
		_, file, line, _ = runtime.Caller(2)
	}

	logger := log.New(os.Stderr, "\n\n\x1b[31m", log.LstdFlags)
	logger.Printf("[Error]:\n%s\n%s%s", msg, stackTrace, "\033[0m")

	//nolint //G401: Use of weak cryptographic primitive
	errorKeyBinary := md5.Sum([]byte(e.appService.Name + ":" + file + ":" + fmt.Sprint(line)))
	errorKey := hex.EncodeToString(errorKeyBinary[:])
	value := &ErrorMessage{
		File:    file,
		Line:    line,
		AppName: e.appService.Name,
		Message: msg,
		Stack:   stackTrace,
	}

	if c != nil {
		c.Request.Body = io.NopCloser(bytes.NewReader(c.Request.Context().Value(e.requestBodyKey).([]byte)))

		requestID, has := c.Get(requestlogger.ID)
		if has {
			value.Request = []byte("X-Request-ID: " + fmt.Sprint(requestID) + "\n\n")
		}

		binaryRequest, _ := httputil.DumpRequest(c.Request, true)
		if len(binaryRequest)*4 <= 64000 {
			value.Request = append(value.Request, binaryRequest...)
		} else {
			value.Request = append(value.Request, []byte("Partial BODY \n\n")...)
			value.Request = append(value.Request, binaryRequest[0:16000]...)
		}
	}

	marshalValue, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	e.redisStorage.HSet(group, errorKey, marshalValue)
	e.redisStorage.HSet(group, errorKey+":time", time.Now().Unix())
	counter := e.redisStorage.HIncrBy(group, errorKey+":counter", 1)

	if group == GroupError &&
		(e.slackService != nil && !e.appService.IsInLocalMode() && !e.appService.IsInTestMode()) &&
		shouldSendNotification(counter) {
		_ = e.slackService.SendToChannel(
			"errors",
			e.slackService.GetErrorChannel(),
			value.Message,
			slackgo.MsgOptionAttachments(
				slackgo.Attachment{
					AuthorName: e.appService.Name,
					Title:      "Error link",
					TitleLink:  e.slackService.GetDevPanelURL() + "#err-" + errorKey,
					Text:       "Counter: " + fmt.Sprint(counter) + " ENV: " + e.appService.Mode,
				},
			),
		)
	}

	if (e.sentryService != nil && !e.appService.IsInLocalMode() && !e.appService.IsInTestMode() && !e.appService.IsInQAMode()) &&
		shouldSendNotification(counter) {
		e.sentryService.CaptureException(fmt.Errorf(value.Message))
	}
}

func shouldSendNotification(counter int64) bool {
	if counter >= 1 && counter <= 20 {
		return true
	}

	if counter == 50 || counter == 100 {
		return true
	}

	if counter < 1000 {
		return false
	}

	for counter > 1 {
		if counter%10 != 0 {
			return false
		}

		counter /= 10
	}

	return true
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)

		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}

	return buf.Bytes()
}

func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}

	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}

	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}

	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}

	name = bytes.Replace(name, centerDot, dot, -1)

	return name
}
