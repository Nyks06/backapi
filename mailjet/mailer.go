package mailjet

import (
	"context"

	"github.com/mailjet/mailjet-apiv3-go"
	webcore "github.com/nyks06/backapi"
)

type Mailer struct {
	FromEmail string
	FromName  string
	Client    *mailjet.Client
}

func NewMailer(fromEmail, fromName, apiKey, apiSecret string) *Mailer {
	return &Mailer{
		FromEmail: fromEmail,
		FromName:  fromName,
		Client:    mailjet.NewMailjetClient(apiKey, apiSecret),
	}
}

func (m *Mailer) Send(ctx context.Context, Name string, Phone string, Email string, Message string) error {
	param := &mailjet.InfoSendMail{
		FromEmail: m.FromEmail,
		FromName:  m.FromName,
		Recipients: []mailjet.Recipient{
			mailjet.Recipient{
				Email: Email,
			},
		},
		Subject:  "Demande de support",
		HTMLPart: Message,
	}
	_, err := m.Client.SendMail(param)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}
