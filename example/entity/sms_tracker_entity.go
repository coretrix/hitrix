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
	ID                    uint64
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

func (s *SmsTrackerEntity) SetStatus(status string) {
	s.Status = status
}

func (s *SmsTrackerEntity) SetTo(to string) {
	s.To = to
}

func (s *SmsTrackerEntity) SetText(text string) {
	s.Text = text
}

func (s *SmsTrackerEntity) SetFromPrimaryGateway(primary string) {
	s.FromPrimaryGateway = primary
}

func (s *SmsTrackerEntity) SetFromSecondaryGateway(secondary string) {
	s.FromSecondaryGateway = secondary
}

func (s *SmsTrackerEntity) SetPrimaryGatewayError(primaryError string) {
	s.PrimaryGatewayError = primaryError
}

func (s *SmsTrackerEntity) SetSecondaryGatewayError(secondaryError string) {
	s.SecondaryGatewayError = secondaryError
}

func (s *SmsTrackerEntity) SetType(typ string) {
	s.Type = typ
}

func (s *SmsTrackerEntity) SetSentAt(sendAt time.Time) {
	s.SentAt = sendAt
}
