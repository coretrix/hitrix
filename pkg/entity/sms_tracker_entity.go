package entity

import (
	"time"

	"github.com/latolukasz/orm"
)

const (
	SMSTrackerTypeSMS     = "sms"
	SMSTrackerTypeCallout = "callout"
)

type smsTrackerTypeAll struct {
	orm.EnumModel
	SMSTrackerTypeSMS     string
	SMSTrackerTypeCallout string
}

var SMSTrackerTypeAll = &smsTrackerTypeAll{
	SMSTrackerTypeSMS:     SMSTrackerTypeSMS,
	SMSTrackerTypeCallout: SMSTrackerTypeCallout,
}

type SmsTrackerEntity struct {
	orm.ORM               `orm:"table=sms_tracker"`
	ID                    uint
	Status                string
	To                    string `orm:"varchar=15"`
	Text                  string
	FromPrimaryGateway    string
	FromSecondaryGateway  string
	PrimaryGatewayError   string
	SecondaryGatewayError string
	Type                  string    `orm:"enum=entity.SMSTrackerTypeAll;required"`
	SentAt                time.Time `orm:"time"`
}
