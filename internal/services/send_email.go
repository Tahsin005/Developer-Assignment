package services

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/tahsin005/affpilot-auth/internal/config"
)

func SendEmail(email, textToBeSent, subject string) {
    cfg := config.GetConfig()
	from := cfg.EmailCfg.From
    pass := cfg.EmailCfg.Password
    to := email

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject:" + subject + "\n\n" +
        textToBeSent

    err := smtp.SendMail(fmt.Sprintf("%s:%s", cfg.EmailCfg.Host, cfg.EmailCfg.Port),
        smtp.PlainAuth("", from, pass, cfg.EmailCfg.Host),
        from, []string{to}, []byte(msg))

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }
    log.Println("Successfully sended to " + to)
}