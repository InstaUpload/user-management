package mail

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

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
	// TODO: Get HTML template for welcome email.
	// Implement the logic to send a welcome email using SMTP
	var rendered bytes.Buffer
	tmpl, err := template.ParseFiles("mail/templates/welcome.html")
	if err != nil {
		log.Printf("Failed to parse template: %s", err.Error())
		return
	}
	if err := tmpl.Execute(&rendered, user); err != nil {
		log.Printf("Failed to execute template: %s", err.Error())
		return
	}
	// Use the rendered HTML content in the email body
	header := "MIME-version: 1.0;\n" + "Content-Type: text/html; charset=\"UTF-8\";\n"
	message := fmt.Sprintf("Subject: Welcome to InstaUpload!\n%s\n\n%s", header, rendered.String())
	host := utils.GetEnvString("MAILHOST", "smtp.example.com")
	post := utils.GetEnvInt("MAILPORT", 587)
	err = smtp.SendMail(
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

func (g *GMailSender) SendVerification(data *types.SendVerificationKM) {
	// TODO: Get HTML template for verification email.
	// TODO: pass on variable like token user's name.
	urlToken := fmt.Sprintf("https://InstaUpload.com/user/verify?token=%s", data.Token)
	tempVariable := struct {
		UrlToken string
	}{
		UrlToken: urlToken,
	}
	var rendered bytes.Buffer
	tmpl, err := template.ParseFiles("mail/templates/verification.html")
	if err != nil {
		log.Printf("Failed to parse template: %s", err.Error())
		return
	}
	if err := tmpl.Execute(&rendered, tempVariable); err != nil {
		log.Printf("Failed to execute template: %s", err.Error())
		return
	}
	// Use the rendered HTML content in the email body
	header := "MIME-version: 1.0;\n" + "Content-Type: text/html; charset=\"UTF-8\";\n"
	message := fmt.Sprintf("Subject: User Varification mail by InstaUpload.\n%s\n\n%s", header, rendered.String())
	host := utils.GetEnvString("MAILHOST", "smtp.example.com")
	post := utils.GetEnvInt("MAILPORT", 587)
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", host, post),
		g.auth,
		utils.GetEnvString("MAILSENDEREMAIL", "gpt.sahaj28@gmail.com"),
		[]string{data.Email},
		[]byte(message),
	)
	if err != nil {
		log.Printf("Failed to send email: %s", err.Error())
		return
	}
}
