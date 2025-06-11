package store_event

import (
	"context"
)

// MySQLStore MYSQL的StoreEvent的实现
type MySQLStore struct {
}

func (m *MySQLStore) Handler(ctx context.Context, data StoreEventInfo) error {
	return nil
}
