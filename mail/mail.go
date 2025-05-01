package mail

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
)

type GMailSender struct {
	auth smtp.Auth
}

func NewMailSender(config types.MailConfig) *GMailSender {
	return &GMailSender{
		auth: smtp.PlainAuth("", config.SenderEmail, config.Password, config.Host),
	}
}

func (g *GMailSender) SendWelcome(user *types.SendWelcomeEmailKM) {
	// Get HTML template for welcome email.
	// Implement the logic to send a welcome email using SMTP
	message := fmt.Sprintf("Subject: Welcome to InstaUpload!\n\nHello %s,\n\nWelcome to InstaUpload! We're glad to have you on board.\n\nBest regards,\nInstaUpload Team", user.Name)
	host := utils.GetEnvString("MAILHOST", "smtp.example.com")
	post := utils.GetEnvInt("MAILPORT", 587)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", host, post),
		g.auth,
		utils.GetEnvString("MAILSENDEREMAIL", "gpt.sahaj28@gmail.com"),
		[]string{user.Email},
		[]byte(message),
	)
	if err != nil {
		log.Printf("Failed to send email: %s", err.Error())
		return
	}
}
