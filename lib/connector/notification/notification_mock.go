package notification

type mockNotification struct {
}

func newMockNotification() mockNotification {
	return mockNotification{}
}

func (mock mockNotification) Send(string, string) error {
	return nil
}
