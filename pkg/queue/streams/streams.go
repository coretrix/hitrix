package streams

const StreamMsgRetryOTP = "msg.retry-otp"

func GetGroupName(queueName string, suffix *string) string {
	groupName := queueName + "_group"
	if suffix != nil {
		groupName += "_" + *suffix
	}
	return groupName
}
