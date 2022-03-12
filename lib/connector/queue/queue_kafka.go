package queue

// load kafka creds to ENV
type kafka struct {
	//client kf.Producer
}

func newKafka() kafka {
	// TODO: create new producer
	return kafka{}
}

func (queue kafka) Send(string, string, string) error {
	// TODO: implement kafka producer
	return nil
}

func (queue kafka) Delete(string, string) error {
	// TODO: implement kafka delete
	return nil
}

func (queue kafka) ReceiveBatch(int, string) ([]Message, error) {
	// TODO: implement kafka consumer
	return []Message{
		Message{},
		Message{},
	}, nil
}
