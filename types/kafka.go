package types

// NOTE: KT stands for Kafka Topic.
// NOTE: KM stands for Kafka Message.

var SendVerificationKT string = "VerifyUser"

type SendVerificationKM struct {
	Token string `json:"token"`
}

type SendWelcomeEmailKM struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
