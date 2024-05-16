package repository

import (
	"context"
	"webook/webook/internal/domain"
)

type HistoryRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
