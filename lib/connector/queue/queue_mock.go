package queue

type mockQueue struct {
}

func newMockQueue() mockQueue {
	return mockQueue{}
}

func (mock mockQueue) Send(string, string, string) error {
	return nil
}

func (mock mockQueue) Delete(string, string) error {
	return nil
}

func (mock mockQueue) ReceiveBatch(int, string) ([]Message, error) {
	return []Message{
		Message{},
		Message{},
	}, nil
}
