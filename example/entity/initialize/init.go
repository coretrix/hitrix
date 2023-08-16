package initialize

import (
	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"
	"github.com/latolukasz/beeorm/v2/plugins/crud_stream"
	"github.com/latolukasz/beeorm/v2/plugins/fake_delete"

	entity2 "github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/redis"
	redissearch "github.com/coretrix/hitrix/example/redis_search"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/queue/streams"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&entity2.APILogEntity{},
		&entity2.AdminUserEntity{},
		&entity2.DevPanelUserEntity{},
		&entity.FileEntity{},
		&entity.SmsTrackerEntity{},
		&entity.OTPTrackerEntity{},
		&entity.FeatureFlagEntity{},
		&entity.RequestLoggerEntity{},
		&entity.RoleEntity{},
		&entity.ResourceEntity{},
		&entity.PrivilegeEntity{},
		&entity.PermissionEntity{},
	)

	registry.RegisterEnumStruct("entity.FileStatusAll", entity.FileStatusAll)
	registry.RegisterEnumStruct("entity.APILogTypeAll", entity2.APILogTypeAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", entity2.APILogStatusAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", entity2.APILogStatusAll)
	registry.RegisterEnumStruct("entity.SMSTrackerTypeAll", entity.SMSTrackerTypeAll)
	registry.RegisterEnumStruct("entity.OTPTrackerTypeAll", entity.OTPTrackerTypeAll)
	registry.RegisterEnumStruct("entity.OTPTrackerGatewaySendStatusAll", entity.OTPTrackerGatewaySendStatusAll)
	registry.RegisterEnumStruct("entity.OTPTrackerGatewayVerifyStatusAll", entity.OTPTrackerGatewayVerifyStatusAll)

	registry.RegisterPlugin(crud_stream.Init(nil))
	registry.RegisterPlugin(fake_delete.Init(nil))

	redisearchPlugin := redisearch.Init(redis.SearchPool)
	redisearchPlugin.RegisterCustomIndex(redissearch.GetDevPanelUserIndex(redis.SearchPool))

	registry.RegisterPlugin(redisearchPlugin)

	// crud_stream plugin stream consumer group
	registry.RegisterRedisStreamConsumerGroups(crud_stream.ChannelName, streams.GetGroupName(crud_stream.ChannelName, nil))

	// redis search indexer
	registry.RegisterRedisStream(redisearch.RedisSearchIndexerChannel, redis.DefaultPool)
	registry.RegisterRedisStreamConsumerGroups(redisearch.RedisSearchIndexerChannel, streams.GetGroupName(redisearch.RedisSearchIndexerChannel, nil))
}
