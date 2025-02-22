package services

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"ticket-processor/internal/ierrors"
	"ticket-processor/internal/models"
	"ticket-processor/internal/receipt"
	"ticket-processor/internal/storage"
	"time"
)

type ReceiptProcessor interface {
	ProcessReceipt(ctx context.Context, receipt models.Receipt) (string, error)
	GetPoints(ctx context.Context, id string) (int, error)
}

type receiptProcessor struct {
	storage storage.Storage
	cache   storage.Cache
	log     *zap.Logger
}

func NewReceiptProcessor(l *zap.Logger, s storage.Storage, c storage.Cache) ReceiptProcessor {
	return &receiptProcessor{
		storage: s,
		cache:   c,
		log:     l,
	}
}

func (rp *receiptProcessor) ProcessReceipt(ctx context.Context, r models.Receipt) (string, error) {
	points := receipt.CalculatePoints(r)

	id, err := rp.storage.Store(ctx, points)
	if err != nil {
		rp.log.Error("Error storing processed receipt", zap.Error(err))
		return "", fmt.Errorf("error storing receipt: %w", err)
	}

	go rp.updateCache(id, points)

	return id, nil
}

func (rp *receiptProcessor) GetPoints(ctx context.Context, id string) (int, error) {
	points, ok := rp.cache.Load(ctx, id)
	if ok {
		return points, nil
	}

	points, ok = rp.storage.Retrieve(ctx, id)
	if !ok {
		return 0, ierrors.ErrNotFound
	}

	go rp.updateCache(id, points)
	return points, nil
}

func (rp *receiptProcessor) updateCache(id string, points int) {
	cacheCtx := context.Background()
	err := rp.cache.Set(cacheCtx, id, points, time.Minute*5)
	if err != nil {
		rp.log.Error("Error setting cache", zap.Error(err))
	}
}
