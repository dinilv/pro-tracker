package bucket

type bucketMock struct {
}

func newBucketMock() bucketMock {
	return bucketMock{}
}

func (mock bucketMock) Upload(string, string, string, []byte) error {
	return nil
}
