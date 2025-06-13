package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

const (
	SMSTrackerTypeSMS     = "sms"
	SMSTrackerTypeCallout = "callout"
)

type smsTrackerTypeAll struct {
	SMSTrackerTypeSMS     string
	SMSTrackerTypeCallout string
}

var SMSTrackerTypeAll = smsTrackerTypeAll{
	SMSTrackerTypeSMS:     SMSTrackerTypeSMS,
	SMSTrackerTypeCallout: SMSTrackerTypeCallout,
}

type SmsTrackerEntity struct {
	beeorm.ORM            `orm:"table=sms_tracker"`
	ID                    uint64
	Status                string
	To                    string `orm:"length=15"`
	Text                  string `orm:"length=max"`
	FromPrimaryGateway    string
	FromSecondaryGateway  string
	PrimaryGatewayError   string    `orm:"length=max"`
	SecondaryGatewayError string    `orm:"length=max"`
	Type                  string    `orm:"enum=entity.SMSTrackerTypeAll;required"`
	SentAt                time.Time `orm:"time"`
}
