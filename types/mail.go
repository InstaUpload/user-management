package types

const MailWelcomeKey = "WelcomeMail"
const MailVerificationKey = "VerificationMail"
const MailEditorInviteKey = "MailEditorInvite"

type MailConfig struct {
	Host        string
	Port        int
	SenderEmail string
	Password    string
}
