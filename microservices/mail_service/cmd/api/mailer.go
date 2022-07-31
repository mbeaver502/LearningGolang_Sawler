package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

var templatesPath = "./templates/"

// Mail is the model for the mail server.
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

// Message represents content for an individual email.
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// SendSMTPMessage sends an email via SMTP.
func (m *Mail) SendSMTPMessage(msg Message) error {
	// use defaults for From fields if none specified
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// create data map that can be passed to templates
	// "message" is the name we'll reference in the template
	data := map[string]any{
		"message": msg.Data,
	}
	msg.DataMap = data

	// create an HTML version of the email
	formattedMsg, err := m.buildHTMLMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	// create a plaintext version of the email
	plainMsg, err := m.buildPlainTextMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	// set up the SMTP server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// connect to our SMTP server
	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	// set up our email message
	email := mail.NewMSG().
		SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetBody(mail.TextPlain, plainMsg).
		AddAlternative(mail.TextHTML, formattedMsg)

	// add attachments, if any, to the email
	if len(msg.Attachments) > 0 {
		for _, a := range msg.Attachments {
			email.AddAttachment(a) // this is deprecated, but it doesn't really matter since we're not using this
		}
	}

	// actually try to send the email
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s%s", templatesPath, "mail.html.gohtml")

	// try to read in the template
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// attempt to execute our template with our data
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// inline the CSS in our template
	formattedMsg := tpl.String()
	formattedMsg, err = m.inlineCSS(formattedMsg)
	if err != nil {
		return "", err
	}

	return formattedMsg, nil
}

func (m *Mail) inlineCSS(msg string) (string, error) {
	opts := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(msg, &opts)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s%s", templatesPath, "mail.plain.gohtml")

	// try to read in the template
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// attempt to execute our template with our data
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMsg := tpl.String()

	return plainMsg, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
