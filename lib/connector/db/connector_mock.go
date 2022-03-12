package db

type mockDB struct {
}

func newMockDB() mockDB {
	return mockDB{}
}

func (db mockDB) Save(string, interface{}) error {
	return nil
}

func (db mockDB) List(int, int, string, []string, map[string]interface{}, map[string]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
