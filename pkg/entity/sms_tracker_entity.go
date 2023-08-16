package entity

import (
	"time"

	"github.com/latolukasz/beeorm/v2"
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

func (s *SmsTrackerEntity) SetStatus(status string) {
	s.Status = status
}

func (s *SmsTrackerEntity) SetTo(to string) {
	s.To = to
}

func (s *SmsTrackerEntity) SetText(text string) {
	s.Text = text
}

func (s *SmsTrackerEntity) SetFromPrimaryProvider(primary string) {
	s.FromPrimaryGateway = primary
}

func (s *SmsTrackerEntity) SetFromSecondaryProvider(secondary string) {
	s.FromSecondaryGateway = secondary
}

func (s *SmsTrackerEntity) SetPrimaryProviderError(primaryError string) {
	s.PrimaryGatewayError = primaryError
}

func (s *SmsTrackerEntity) SetSecondaryProviderError(secondaryError string) {
	s.SecondaryGatewayError = secondaryError
}

func (s *SmsTrackerEntity) SetType(typ string) {
	s.Type = typ
}

func (s *SmsTrackerEntity) SetSentAt(sendAt time.Time) {
	s.SentAt = sendAt
}
