package entity

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/latolukasz/beeorm"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&AdminUserEntity{}, &APILogEntity{}, &entity.SmsTrackerEntity{},
	)

	registry.RegisterEnumStruct("entity.APILogTypeAll", APILogTypeAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.SMSTrackerTypeAll", entity.SMSTrackerTypeAll)
}
