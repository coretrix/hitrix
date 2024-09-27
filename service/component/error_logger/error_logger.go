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
	"math"
	"net/http/httputil"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"
	slackgo "github.com/slack-go/slack"

	"github.com/coretrix/hitrix/service/component/app"
	requestlogger "github.com/coretrix/hitrix/service/component/request_logger"
	"github.com/coretrix/hitrix/service/component/sentry"
	"github.com/coretrix/hitrix/service/component/slack"
)

const GroupError = "error"

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
	return &RedisErrorLogger{
		redisStorage:   ormService.GetRedis(),
		slackService:   slackService,
		appService:     appService,
		sentryService:  sentryService,
		requestBodyKey: requestBodyKey,
	}
}

func (e *RedisErrorLogger) LogError(errData interface{}) {
	e.log(errData, 2, nil)
}

func (e *RedisErrorLogger) LogErrorWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, 2, c)
}

func (e *RedisErrorLogger) LogPanicWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, 4, c)
}

func (e *RedisErrorLogger) log(errData interface{}, callerSkip int, c *gin.Context) {
	var msg string

	err, ok := errData.(error)
	if ok {
		msg = err.Error()
	} else {
		msg = errData.(string)
	}

	logger := log.New(os.Stderr, "\n\n\x1b[31m", log.LstdFlags)
	stackTrace := stack(0)
	logger.Printf("[Error]:\n%s\n%s%s", msg, stackTrace, "\033[0m")

	_, file, line, _ := runtime.Caller(callerSkip)

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

	e.redisStorage.HSet(GroupError, errorKey, marshalValue)
	e.redisStorage.HSet(GroupError, errorKey+":time", time.Now().Unix())
	counter := e.redisStorage.HIncrBy(GroupError, errorKey+":counter", 1)

	logg := math.Log10(float64(counter))

	if (e.slackService != nil && !e.appService.IsInLocalMode() && !e.appService.IsInTestMode()) && logg == float64(int64(logg)) {
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
		logg == float64(int64(logg)) {
		e.sentryService.CaptureException(fmt.Errorf(value.Message))
	}
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
