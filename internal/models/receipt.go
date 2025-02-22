package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Item struct {
	ShortDescription string  `json:"shortDescription" validate:"required,notblank"`
	Price            float64 `json:"price" validate:"required"`
}

type Receipt struct {
	Retailer     string  `json:"retailer" validate:"required,notblank"`
	PurchaseDate string  `json:"purchaseDate" validate:"required,date"`
	PurchaseTime string  `json:"purchaseTime" validate:"required,time"`
	Items        []Item  `json:"items" validate:"required,dive"`
	Total        float64 `json:"total" validate:"required"`
}
type ProcessReceiptResponse struct {
	ID string `json:"id"`
}
type GetReceiptPointsResponse struct {
	Points int `json:"points"`
}

func (r *Receipt) UnmarshalJSON(data []byte) error {
	var raw struct {
		Retailer     string `json:"retailer"`
		PurchaseDate string `json:"purchaseDate"`
		PurchaseTime string `json:"purchaseTime"`
		Items        []Item `json:"items"`
		Total        string `json:"total"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Retailer = raw.Retailer
	r.PurchaseDate = raw.PurchaseDate
	r.PurchaseTime = raw.PurchaseTime
	r.Items = raw.Items

	total, err := strconv.ParseFloat(raw.Total, 64)
	if err != nil {
		return fmt.Errorf("invalid total format: %w", err)
	}
	r.Total = total

	return nil
}

func (i *Item) UnmarshalJSON(data []byte) error {
	var raw struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	i.ShortDescription = raw.ShortDescription

	price, err := strconv.ParseFloat(raw.Price, 64)
	if err != nil {
		return fmt.Errorf("invalid price format: %w", err)
	}
	i.Price = price

	return nil
}
