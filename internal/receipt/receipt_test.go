package receipt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"ticket-processor/internal/models"
)

func TestCalculatePoints_ValidReceipt(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int
	}{
		{
			name: "Example1",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: 6.49},
					{ShortDescription: "Emils Cheese Pizza", Price: 12.25},
					{ShortDescription: "Knorr Creamy Chicken", Price: 1.26},
					{ShortDescription: "Doritos Nacho Cheese", Price: 3.35},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: 12.00},
				},
				Total: 35.35,
			},
			expected: 28,
		},
		{
			name: "Example2",
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: 2.25},
					{ShortDescription: "Gatorade", Price: 2.25},
					{ShortDescription: "Gatorade", Price: 2.25},
					{ShortDescription: "Gatorade", Price: 2.25},
				},
				Total: 9.00,
			},
			expected: 109,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := CalculatePoints(tt.receipt)
			assert.Equal(t, tt.expected, points)
		})
	}
}

func TestCalculateRetailerPoints(t *testing.T) {
	tests := []struct {
		name     string
		retailer string
		expected int
	}{
		{
			name:     "EmptyRetailer",
			retailer: "",
			expected: 0,
		},
		{
			name:     "SimpleRetailer",
			retailer: "Target",
			expected: 6,
		},
		{
			name:     "RetailerWithNumbers",
			retailer: "Target123",
			expected: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := calculateRetailerPoints(tt.retailer)
			assert.Equal(t, tt.expected, points)
		})
	}
}

func TestCalculateQuarterMultipleBonus(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected int
	}{
		{
			name:     "QuarterMultiple",
			amount:   10.25,
			expected: 25,
		},
		{
			name:     "NotQuarterMultiple",
			amount:   10.60,
			expected: 0,
		},
		{
			name:     "RoundDollar",
			amount:   10.00,
			expected: 25,
		},
		{
			name:     "ZeroAmount",
			amount:   0.00,
			expected: 0,
		},
		{
			name:     "NegativeAmount",
			amount:   -10.25,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := calculateQuarterMultipleBonus(tt.amount)
			assert.Equal(t, tt.expected, points)
		})
	}
}

func TestCalculateRoundDollarBonus(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected int
	}{
		{
			name:     "RoundDollar",
			amount:   10.00,
			expected: 50,
		},
		{
			name:     "NonRoundDollar",
			amount:   10.50,
			expected: 0,
		},
		{
			name:     "ZeroAmount",
			amount:   0.00,
			expected: 0,
		},
		{
			name:     "NegativeAmount",
			amount:   -10.00,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := calculateRoundDollarBonus(tt.amount)
			assert.Equal(t, tt.expected, points)
		})
	}
}

func TestCalculateItemCountBonus(t *testing.T) {
	tests := []struct {
		name      string
		itemCount int
		expected  int
	}{
		{
			name:      "EvenItemCount",
			itemCount: 4,
			expected:  10,
		},
		{
			name:      "OddItemCount",
			itemCount: 5,
			expected:  10,
		},
		{
			name:      "ZeroItemCount",
			itemCount: 0,
			expected:  0,
		},
		{
			name:      "NegativeItemCount",
			itemCount: -1,
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := calculateItemCountBonus(tt.itemCount)
			assert.Equal(t, tt.expected, points)
		})
	}
}

func TestCalculateItemDescriptionPoints(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.Item
		expected int
	}{
		{
			name:     "EmptyItems",
			items:    []models.Item{},
			expected: 0,
		},
		{
			name: "MultipleItems",
			items: []models.Item{
				{ShortDescription: "Candy Bar", Price: 0.75},
				{ShortDescription: "2 lb Chicken", Price: 10.00},
				{ShortDescription: "Milk", Price: 4.50},
			},
			expected: 3,
		},
		{
			name: "ItemWithDescriptionMultipleOfThreeAndRoundUp",
			items: []models.Item{
				{ShortDescription: "ItemDescription", Price: 10.50},
			},
			expected: 3,
		},
		{
			name: "ItemWithDescriptionMultipleOfThreeWithoutRoundUp",
			items: []models.Item{
				{ShortDescription: "ItemDescription", Price: 10.00},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := calculateItemDescriptionPoints(tt.items)
			assert.Equal(t, tt.expected, points)
		})
	}
}

func TestCalculateOddDayBonus(t *testing.T) {
	tests := []struct {
		name         string
		purchaseDate string
		expected     int
	}{
		{
			name:         "OddDay",
			purchaseDate: "2023-10-27",
			expected:     6,
		},
		{
			name:         "EvenDay",
			purchaseDate: "2023-10-28",
			expected:     0,
		},
		{
			name:         "InvalidDate",
			purchaseDate: "invalid date",
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := calculateOddDayBonus(tt.purchaseDate)
			assert.Equal(t, tt.expected, points)
		})
	}
}

func TestCalculatePurchaseTimeBonus(t *testing.T) {
	tests := []struct {
		name         string
		purchaseTime string
		expected     int
	}{
		{
			name:         "Between2and4PM",
			purchaseTime: "15:00",
			expected:     10,
		},
		{
			name:         "Before2PM",
			purchaseTime: "13:00",
			expected:     0,
		},
		{
			name:         "After4PM",
			purchaseTime: "17:00",
			expected:     0,
		},
		{
			name:         "InvalidTime",
			purchaseTime: "invalid time",
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := calculatePurchaseTimeBonus(tt.purchaseTime)
			assert.Equal(t, tt.expected, points)
		})
	}
}
