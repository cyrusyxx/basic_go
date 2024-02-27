package sms

import "context"

type Service interface {
	// Send sends a SMS to the given number with the given template ID and arguments.
	Send(ctx context.Context, number string,
		tplID string, args []string, numbers ...string) error
}
