package entity

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/latolukasz/beeorm"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&DevPanelUserEntity{}, &APILogEntity{}, &entity.SmsTrackerEntity{}, &entity.OTPTrackerEntity{},
	)

	registry.RegisterEnumStruct("entity.APILogTypeAll", APILogTypeAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.SMSTrackerTypeAll", entity.SMSTrackerTypeAll)
	registry.RegisterEnumStruct("entity.OTPTrackerTypeAll", entity.OTPTrackerTypeAll)
	registry.RegisterEnumStruct("entity.OTPTrackerGatewaySendStatusAll", entity.OTPTrackerGatewaySendStatusAll)
	registry.RegisterEnumStruct("entity.OTPTrackerGatewayVerifyStatusAll", entity.OTPTrackerGatewayVerifyStatusAll)
}
