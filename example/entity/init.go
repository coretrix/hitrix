package entity

import (
	"github.com/latolukasz/orm"
)

func Init(registry *orm.Registry) {
	registry.RegisterEntity(
		&AdminUserEntity{}, &APILogEntity{}, &SmsTrackerEntity{},
	)

	registry.RegisterEnumStruct("entity.APILogTypeAll", APILogTypeAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.APILogStatusAll", APILogStatusAll)
	registry.RegisterEnumStruct("entity.SMSTrackerTypeAll", SMSTrackerTypeAll)
}
