package entity

import (
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&DevPanelUserEntity{},
		&APILogEntity{},
		&entity.SmsTrackerEntity{},
		&entity.OTPTrackerEntity{},
		&entity.FeatureFlagEntity{},
		&entity.RequestLoggerEntity{},
	)

	registry.RegisterEnumStruct("entity.APILogTypeAll", APILogTypeAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.SMSTrackerTypeAll", entity.SMSTrackerTypeAll)
	registry.RegisterEnumStruct("entity.OTPTrackerTypeAll", entity.OTPTrackerTypeAll)
	registry.RegisterEnumStruct("entity.OTPTrackerGatewaySendStatusAll", entity.OTPTrackerGatewaySendStatusAll)
	registry.RegisterEnumStruct("entity.OTPTrackerGatewayVerifyStatusAll", entity.OTPTrackerGatewayVerifyStatusAll)
}
