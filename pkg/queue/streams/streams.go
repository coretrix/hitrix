package streams

const StreamMsgRetryOTP = "msg.retry-otp"

func GetGroupName(queueName string, suffix *string) string {
	if suffix == nil {
		return queueName + "_group"
	}

	return queueName + "_group_" + *suffix
}
