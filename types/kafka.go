package types

// NOTE: KT stands for Kafka Topic.
// NOTE: KM stands for Kafka Message.

const EmailUserKT = "EmailUser"

type SendVerificationKM struct {
	Token string `json:"token"`
}

type SendWelcomeEmailKM struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
