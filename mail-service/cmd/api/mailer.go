package main

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

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

type Message struct {
	From       string
	FromName   string
	To         string
	Subject    string
	Attachment []string
	Data       any
	DataMap    map[string]any
}

var tmplFS embed.FS

func (m *Mail) SendSMTPMessage(msg Message) error {

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMsg, err := m.buildHTMLMsg(msg)
	if err != nil{
		return err
	}

	plainMsg, err := m.buildPlainMsg(msg)
	if err != nil{
		return err
	} 


	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Encryption = m.getEncryption(m.Encryption)
	server.Username = m.Username
	server.Password = m.Password
	server.KeepAlive = false
	server.ConnectTimeout = time.Second * 10
	server.SendTimeout = time.Second * 10

	smtpClient, err := server.Connect()
	if err != nil{
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMsg)
	email.AddAlternative(mail.TextHTML, formattedMsg)
	if len(msg.Attachment ) > 0 {
		for _, x := range msg.Attachment{
			 email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil{
		return err
	}


	return nil
}

func (m *Mail) buildHTMLMsg(msg Message) (string, error) {

	templateToRender := "./templates/mail.html.gohtml"
	
	temp, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil{
		return "", err
	}

	var tpl bytes.Buffer

	if err = temp.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil{
		return "", nil
	}

	formattedMsg := tpl.String()
	formattedMsg, err = m.inlineCSS(formattedMsg)
	if err != nil{
		return "", err
	}


	return formattedMsg, nil

}

func (m *Mail) buildPlainMsg(msg Message) (string, error) {

	templateToRender := "./templates/mail.plain.gohtml"

	temp, err := template.New("email-plaintext").ParseFiles(templateToRender)
	if err != nil{
		return "", err
	}

	var tpl bytes.Buffer

	if err = temp.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil{
		return "", nil
	}

	plainMsg := tpl.String()

	return plainMsg, nil

}

func (m *Mail) inlineCSS(s string) (string, error){

	options := premailer.Options{
		RemoveClasses: false,
		CssToAttributes: false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil{
		return "", err
	}

	html, err := prem.Transform()
	if err != nil{
		return "", err
	}

	return html, nil
}


func(m *Mail) getEncryption(s string) mail.Encryption{
	
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSL
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}