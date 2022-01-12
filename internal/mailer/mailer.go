/*
The mailer package is responsible for sending smtp email tempaltes
*/
package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"time"

	"gopkg.in/gomail.v2"
)

//go:embed templates
var templateFS embed.FS

// The Mailer type holds the dailer, sender, and templates
type Mailer struct {
	dialer    *gomail.Dialer
	sender    string
	templates fs.FS
}

func New(host string, port int, username, password, sender string) (Mailer, error) {
	templates, err := fs.Sub(templateFS, "templates")
	if err != nil {
		return Mailer{}, err
	}

	dialer := gomail.NewDialer(host, port, username, password)

	return Mailer{
		dialer:    dialer,
		sender:    sender,
		templates: templates,
	}, nil
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	t, err := template.ParseFS(m.templates, templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = t.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = t.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = t.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// retry 3 times  if error with 500ms delay
	for i := 1; i <= 3; i++ {

		err = m.dialer.DialAndSend(msg)
		if nil == err {
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return err

}
