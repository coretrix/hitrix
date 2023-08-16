package main

import (
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity/initialize"
	"github.com/coretrix/hitrix/example/redis"
	"github.com/coretrix/hitrix/example/scripts"
	scriptsHitrix "github.com/coretrix/hitrix/scripts"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/registry"
)

func main() {
	s, deferFunc := hitrix.New(
		"single-instance-cron", "secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("../../config"),
		registry.ServiceProviderOrmRegistry(initialize.Init),
		registry.ServiceProviderOrmEngine(redis.SearchPool),
		registry.ServiceProviderClock(),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(false, redis.SearchPool),
	).RegisterRedisPools(&app.RedisPools{
		Cache:      redis.DefaultPool,
		Persistent: redis.DefaultPool,
		Stream:     redis.StreamsPool,
		Search:     redis.SearchPool,
	}).Build()
	defer deferFunc()

	b := &hitrix.BackgroundProcessor{Server: s}
	b.RunAsyncOrmConsumer()
	b.RunAsyncRequestLoggerCleaner()

	s.RunBackgroundProcess(func(b *hitrix.BackgroundProcessor) {
		go b.RunScript(&scripts.DirtyConsumerScript{})
		go b.RunScript(&scriptsHitrix.ReindexConsumerScript{})
	})
}
