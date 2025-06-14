package types

// NOTE: KT stands for Kafka Topic.
// NOTE: KM stands for Kafka Message.

const EmailUserKT = "EmailUser"

type SendVerificationKM struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SendWelcomeEmailKM struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SendEditorRequestKM struct {
	CreatorName string `json:"creater_name"`
	EditorEmail string `json:"email"`
	Token       string `json:"token"`
}
