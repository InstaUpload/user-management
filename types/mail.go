package types

const MailWelcomeKey = "WelcomeMail"

type MailConfig struct {
	Host        string
	Port        int
	SenderEmail string
	Password    string
}
