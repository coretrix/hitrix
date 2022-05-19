package consumers

import (
	"log"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/queue/streams"
	"github.com/coretrix/hitrix/service/component/otp"
)

type OTPRetryConsumer struct {
	ormService *beeorm.Engine
	maxRetries int
}

func NewOTPRetryConsumer(ormService *beeorm.Engine, maxRetries int) *OTPRetryConsumer {
	return &OTPRetryConsumer{ormService: ormService, maxRetries: maxRetries}
}

func (c *OTPRetryConsumer) GetQueueName() string {
	return streams.StreamMsgRetryOTP
}

func (c *OTPRetryConsumer) GetGroupName(suffix *string) string {
	return streams.GetGroupName(c.GetQueueName(), suffix)
}

func (c *OTPRetryConsumer) Consume(_ *beeorm.Engine, event beeorm.Event) error {
	log.Println(".")

	ormService := c.ormService.Clone()

	retryDTO := &otp.RetryDTO{}
	event.Unserialize(retryDTO)

	otpTrackerEntity := &entity.OTPTrackerEntity{}
	ormService.LoadByID(retryDTO.OTPTrackerEntityID, otpTrackerEntity)

	RetryOTP(ormService, retryDTO, otpTrackerEntity, c.maxRetries)

	return nil
}

func RetryOTP(ormService *beeorm.Engine, retryDTO *otp.RetryDTO, otpTrackerEntity *entity.OTPTrackerEntity, maxRetries int) {
	retryAfter := time.Second / 2

	retryCount := 1
	for retryCount <= maxRetries {
		retryCount++

		var err error
		otpTrackerEntity.GatewaySendRequest, otpTrackerEntity.GatewaySendResponse, err = retryDTO.Gateway.SendOTP(retryDTO.Phone, retryDTO.Code)
		if err == nil {
			otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusSent
		}

		otpTrackerEntity.RetryCount = retryCount - 1
		if retryCount == maxRetries {
			otpTrackerEntity.MaxRetriesReached = true
		}

		ormService.Flush(otpTrackerEntity)

		if otpTrackerEntity.GatewaySendStatus == entity.OTPTrackerGatewaySendStatusSent {
			break
		}

		time.Sleep(retryAfter)
		retryAfter = retryAfter * 2
	}
}
