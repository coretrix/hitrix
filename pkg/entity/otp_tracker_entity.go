package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

const (
	OTPTrackerTypeSMS   = "sms"
	OTPTrackerTypeEmail = "email"
)

type otpTrackerTypeAll struct {
	OTPTrackerTypeSMS   string
	OTPTrackerTypeEmail string
}

var OTPTrackerTypeAll = otpTrackerTypeAll{
	OTPTrackerTypeSMS:   OTPTrackerTypeSMS,
	OTPTrackerTypeEmail: OTPTrackerTypeEmail,
}

const (
	OTPTrackerGatewaySendStatusNew          = "new"
	OTPTrackerGatewaySendStatusGatewayError = "gateway_error"
	OTPTrackerGatewaySendStatusSent         = "sent"
)

type OTPTrackerGatewaySendStatus struct {
	OTPTrackerGatewaySendStatusNew          string
	OTPTrackerGatewaySendStatusGatewayError string
	OTPTrackerGatewaySendStatusSent         string
}

var OTPTrackerGatewaySendStatusAll = OTPTrackerGatewaySendStatus{
	OTPTrackerGatewaySendStatusNew:          OTPTrackerGatewaySendStatusNew,
	OTPTrackerGatewaySendStatusGatewayError: OTPTrackerGatewaySendStatusGatewayError,
	OTPTrackerGatewaySendStatusSent:         OTPTrackerGatewaySendStatusSent,
}

const (
	OTPTrackerGatewayVerifyStatusNew          = "new"
	OTPTrackerGatewayVerifyStatusGatewayError = "gateway_error"
	OTPTrackerGatewayVerifyStatusInvalidCode  = "invalid_code"
	OTPTrackerGatewayVerifyStatusSuccess      = "success"
	OTPTrackerGatewayVerifyStatusExpired      = "expired"
)

type OTPTrackerGatewayVerifyStatus struct {
	OTPTrackerGatewayVerifyStatusNewError     string
	OTPTrackerGatewayVerifyStatusGatewayError string
	OTPTrackerGatewayVerifyStatusInvalidCode  string
	OTPTrackerGatewayVerifyStatusSuccess      string
	OTPTrackerGatewayVerifyStatusExpired      string
}

var OTPTrackerGatewayVerifyStatusAll = OTPTrackerGatewayVerifyStatus{
	OTPTrackerGatewayVerifyStatusNewError:     OTPTrackerGatewayVerifyStatusNew,
	OTPTrackerGatewayVerifyStatusGatewayError: OTPTrackerGatewayVerifyStatusGatewayError,
	OTPTrackerGatewayVerifyStatusInvalidCode:  OTPTrackerGatewayVerifyStatusInvalidCode,
	OTPTrackerGatewayVerifyStatusSuccess:      OTPTrackerGatewayVerifyStatusSuccess,
	OTPTrackerGatewayVerifyStatusExpired:      OTPTrackerGatewayVerifyStatusExpired,
}

type OTPTrackerEntity struct {
	beeorm.ORM            `orm:"table=otp_tracker"`
	ID                    uint64
	Type                  string `orm:"enum=entity.OTPTrackerTypeAll;required"`
	To                    string `orm:"length=50"`
	Code                  string
	GatewayName           string
	GatewayPriority       uint8
	GatewaySendStatus     string `orm:"enum=entity.OTPTrackerGatewaySendStatusAll;required"`
	GatewaySendRequest    string `orm:"length=max"`
	GatewaySendResponse   string `orm:"length=max"`
	GatewayVerifyStatus   string `orm:"enum=entity.OTPTrackerGatewayVerifyStatusAll;required"`
	GatewayVerifyRequest  string `orm:"length=max"`
	GatewayVerifyResponse string `orm:"length=max"`
	RetryCount            int
	MaxRetriesReached     bool
	SentAt                time.Time `orm:"time"`
}
