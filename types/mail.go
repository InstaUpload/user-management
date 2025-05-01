package types

const MailWelcomeKey = "WelcomeMail"
const MailVerificationKey = "VerificationMail"

type MailConfig struct {
	Host        string
	Port        int
	SenderEmail string
	Password    string
}
