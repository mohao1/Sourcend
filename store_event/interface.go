package store_event

import (
	"context"
)

// StoreEvent StoreEvent的接口
type StoreEvent interface {
	Handler(ctx context.Context, data StoreEventInfo) error
}
