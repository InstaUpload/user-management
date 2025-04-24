package types

// NOTE: KT stands for Kafka Topic.
// NOTE: KM stands for Kafka Message.

var SendVerificationKT string = "VerifyUser"

type SendVerificationKM struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
