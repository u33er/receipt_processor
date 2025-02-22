package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"ticket-processor/internal/models"
	"time"
)

type mockStorage struct {
	storeFunc    func(ctx context.Context, points int) (string, error)
	retrieveFunc func(ctx context.Context, key string) (int, bool)
}

func (m *mockStorage) Store(ctx context.Context, points int) (string, error) {
	return m.storeFunc(ctx, points)
}

func (m *mockStorage) Retrieve(ctx context.Context, key string) (int, bool) {
	return m.retrieveFunc(ctx, key)
}

type mockCache struct {
	getFunc func(ctx context.Context, id string) (int, bool)
	setFunc func(ctx context.Context, id string, points int, ttl time.Duration) error
}

func (m *mockCache) Load(ctx context.Context, id string) (int, bool) {
	return m.getFunc(ctx, id)
}

func (m *mockCache) Set(ctx context.Context, id string, points int, ttl time.Duration) error {
	return m.setFunc(ctx, id, points, ttl)
}

func TestProcessReceipt_Success(t *testing.T) {
	mockStorage := &mockStorage{
		storeFunc: func(ctx context.Context, points int) (string, error) {
			return "receipt-id", nil
		},
	}
	mockCache := &mockCache{
		setFunc: func(ctx context.Context, id string, points int, ttl time.Duration) error {
			return nil
		},
	}
	logger := zap.NewNop()
	rp := NewReceiptProcessor(logger, mockStorage, mockCache)

	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: 6.49},
		},
		Total: 6.49,
	}

	id, err := rp.ProcessReceipt(context.Background(), receipt)
	assert.NoError(t, err)
	assert.Equal(t, "receipt-id", id)
}

func TestProcessReceipt_StorageError(t *testing.T) {
	mockStorage := &mockStorage{
		storeFunc: func(ctx context.Context, points int) (string, error) {
			return "", errors.New("storage error")
		},
	}
	mockCache := &mockCache{
		setFunc: func(ctx context.Context, id string, points int, ttl time.Duration) error {
			return nil
		},
	}
	logger := zap.NewNop()
	rp := NewReceiptProcessor(logger, mockStorage, mockCache)

	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: 6.49},
		},
		Total: 6.49,
	}

	id, err := rp.ProcessReceipt(context.Background(), receipt)
	assert.Error(t, err)
	assert.Empty(t, id)
}
