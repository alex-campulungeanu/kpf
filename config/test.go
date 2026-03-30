package config

import "fmt"

type Notifier interface {
	Send(message string) error
}

type AlertService struct {
	Provider Notifier
}

func (s AlertService) Notify(msg string) error {
	return s.Provider.Send(msg)
}

type EmailProvider struct{}

func (e EmailProvider) Send(msg string) error {
	fmt.Println("Calling the real email provider")
	fmt.Println("Sending email with content", msg)
	return nil
}

func ExecutAlert() {
	service := AlertService{
		Provider: EmailProvider{},
	}
	service.Notify("kind reminder from email")
}
