package webcore

import "context"

type Mailer interface {
	Send(ctx context.Context, Name string, Phone string, Email string, Message string) error
}
