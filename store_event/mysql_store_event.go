package store_event

import (
	"context"
)

type MySQLStore struct {
}

func (m *MySQLStore) Handler(ctx context.Context, data StoreEventInfo) error {
	return nil
}
