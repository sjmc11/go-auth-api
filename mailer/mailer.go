package mailer

import (
	"errors"
	"fmt"
	"github.com/ainsleyclark/go-mail/drivers"
	mail "github.com/ainsleyclark/go-mail/mail"
)

var MailDriver mail.Mailer

func InitGoMailer() (mail.Mailer, error) {
	cfg := mail.Config{
		URL:         "https://api.eu.sparkpost.com",
		APIKey:      "a529be82177a927ec4529b5b83b00be8fe514a7d",
		FromName:    "Resource Link Building",
		FromAddress: "account@lacuna.reddico.io",
	}

	sparkpostDriver, err := drivers.NewSparkPost(cfg)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("something went wrong")
	}

	MailDriver = sparkpostDriver

	return sparkpostDriver, nil
}

//func DoMail(recipient string, subject string, htmlBody string, textBody string){
//
//	tx := &mail.Transmission{
//		Recipients: []string{recipient},
//		Subject:    subject,
//		HTML:       htmlBody,
//		PlainText:  textBody,
//	}
//
//
//	result, err := mailDriver.Send(tx)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(result)
//
//}
