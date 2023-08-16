package scripts

import (
	"context"

	"github.com/coretrix/hitrix/example/redis_stream/consumers"
	"github.com/coretrix/hitrix/pkg/queue"
	"github.com/coretrix/hitrix/service/component/app"
)

type DirtyConsumerScript struct {
}

func (script *DirtyConsumerScript) Run(ctx context.Context, _ app.IExit) {
	queue.NewConsumerRunner(ctx).RunConsumerMany(consumers.NewDirtyConsumer(), nil, 1000)
}

func (script *DirtyConsumerScript) Infinity() bool {
	return true
}

func (script *DirtyConsumerScript) Unique() bool {
	return true
}

func (script *DirtyConsumerScript) Description() string {
	return "dirty consumer script"
}
