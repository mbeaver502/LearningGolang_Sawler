package main

import (
	"time"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()

	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	email := mail.NewMSG()
	email = email.AddTo(m.To).SetFrom(m.From).SetSubject(m.Subject).SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	app.InfoLog.Printf("email sent from %s to %s with subject %s", m.From, m.To, m.Subject)
}
