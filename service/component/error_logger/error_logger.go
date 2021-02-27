package errorlogger

import (
	"bytes"
	"io/ioutil"
	"math"
	"net/http/httputil"

	slackapi "github.com/coretrix/hitrix/service/component/slack_api"

	"github.com/coretrix/hitrix/service/component/app"

	"github.com/bluele/slack"

	//nolint
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/summer-solutions/orm"
)

const GroupError = "error"

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type ErrorLogger interface {
	LogErrorWithRequest(c *gin.Context, errData interface{})
	LogError(dataFromRecover interface{})
	log(errData interface{}, request *http.Request)
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
	redisStorage *orm.RedisCache
	slackService *slackapi.SlackAPI
	appService   *app.App
}

func NewRedisErrorLogger(appService *app.App, ormService *orm.Engine, slackService *slackapi.SlackAPI) ErrorLogger {
	return &RedisErrorLogger{redisStorage: ormService.GetRedis(), slackService: slackService, appService: appService}
}

func (e *RedisErrorLogger) LogError(errData interface{}) {
	e.log(errData, nil)
}

func (e *RedisErrorLogger) LogErrorWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, c.Request)
}

func (e *RedisErrorLogger) log(errData interface{}, request *http.Request) {
	var msg string
	err, ok := errData.(error)
	if ok {
		msg = err.Error()
	} else {
		msg = errData.(string)
	}

	logger := log.New(os.Stderr, "\n\n\x1b[31m", log.LstdFlags)
	stack := stack(0)
	logger.Printf("[Error]:\n%s\n%s%s", msg, stack, "\033[0m")

	_, file, line, _ := runtime.Caller(4)
	//nolint
	errorKeyBinary := md5.Sum([]byte(file + ":" + fmt.Sprint(line)))
	errorKey := hex.EncodeToString(errorKeyBinary[:])
	value := &ErrorMessage{
		File:    file,
		Line:    line,
		AppName: e.appService.Name,
		Message: msg,
		Stack:   stack,
	}

	if request != nil {
		binaryRequest, _ := httputil.DumpRequest(request, true)
		if len(binaryRequest) > 2000 {
			value.Request = binaryRequest[0:2000]
		} else {
			value.Request = binaryRequest
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

	if (e.slackService != nil && e.appService.IsInProdMode() || e.appService.IsInDemoMode()) && logg == float64(int64(logg)) {
		e.slackService.SendToChannel(
			e.slackService.GetErrorChannel(),
			value.Message,
			&slack.ChatPostMessageOpt{
				Attachments: []*slack.Attachment{
					{
						AuthorName: e.appService.Name,
						Title:      "Error link",
						TitleLink:  e.slackService.GetDevPanelURL() + "#err-" + errorKey,
						Text:       "Counter: " + fmt.Sprint(counter) + " ENV: " + e.appService.Mode,
					},
				},
			})
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
			data, err := ioutil.ReadFile(file)
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
