package store

type TReceiptStore map[string]int64

var ReceiptStore TReceiptStore

func CreateReceiptStore() TReceiptStore {
	ReceiptStore = make(map[string]int64)
	return ReceiptStore
}

func (s TReceiptStore) Get(id string) (int64, bool) {
	points, ok := s[id]
	return points, ok
}

func (s TReceiptStore) Set(id string, points int64) {
	s[id] = points
}
