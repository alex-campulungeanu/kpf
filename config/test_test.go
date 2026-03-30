package config

import (
	"fmt"
	"testing"
)

type MockProvider struct {
	lastMessage string
	called      bool
}

func (m *MockProvider) Send(message string) error {
	m.lastMessage = message
	m.called = true
	fmt.Println("calling the mock notifier with message", message)
	return nil
}

func TestAlertService(t *testing.T) {
	mock := &MockProvider{}
	service := AlertService{
		Provider: mock,
	}
	testMessage := "from test_test.go"
	service.Notify(testMessage)
	fmt.Println(mock.lastMessage)
	if !mock.called {
		t.Error("expected Provider.send to be called, but i wasn't")
	}

	if mock.lastMessage != testMessage {
		t.Errorf("expected last message to be %s, but got %s", testMessage, mock.lastMessage)
	}

}
