package consumers

import (
	"fmt"
	"log"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/queue/streams"
	"github.com/coretrix/hitrix/service/component/otp"
)

type OTPRetryConsumer struct {
	ormService      *beeorm.Engine
	maxRetries      int
	gatewayRegistry map[string]otp.IOTPSMSGateway
}

func NewOTPRetryConsumer(ormService *beeorm.Engine, maxRetries int, gatewayRegistry map[string]otp.IOTPSMSGateway) *OTPRetryConsumer {
	return &OTPRetryConsumer{ormService: ormService, maxRetries: maxRetries, gatewayRegistry: gatewayRegistry}
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

	RetryOTP(ormService, c.gatewayRegistry, retryDTO, otpTrackerEntity, c.maxRetries)

	return nil
}

func RetryOTP(ormService *beeorm.Engine, gatewayRegistry map[string]otp.IOTPSMSGateway, retryDTO *otp.RetryDTO, otpTrackerEntity *entity.OTPTrackerEntity, maxRetries int) {
	retryAfter := time.Second / 2

	retryCount := 1
	for retryCount <= maxRetries {
		retryCount++

		gateway, ok := gatewayRegistry[retryDTO.Gateway]
		if !ok {
			panic(fmt.Sprintf("gateway %s not found in registry", gateway))
		}

		var err error
		otpTrackerEntity.GatewaySendRequest, otpTrackerEntity.GatewaySendResponse, err = gateway.SendOTP(retryDTO.Phone, retryDTO.Code)
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
