package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"ticket-processor/internal/ierrors"
	"ticket-processor/internal/models"
	"ticket-processor/internal/services"
	"ticket-processor/internal/validation"
)

type ReceiptHandler interface {
	PostReceiptsProcess(c echo.Context) error
	GetReceiptsIdPoints(c echo.Context) error
}

type receiptHandler struct {
	log              *zap.Logger
	receiptProcessor services.ReceiptProcessor
}

func NewReceiptHandler(log *zap.Logger, receiptProcessor services.ReceiptProcessor) ReceiptHandler {
	return &receiptHandler{
		log:              log,
		receiptProcessor: receiptProcessor,
	}
}

func (h *receiptHandler) PostReceiptsProcess(c echo.Context) error {
	h.log.Info("Processing receipt")

	var receipt models.Receipt

	if err := c.Bind(&receipt); err != nil {
		h.log.Error("Invalid JSON format", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ierrors.NewErrorResponse(http.StatusBadRequest, "Invalid JSON format"))
	}

	if err := validation.ValidateReceipt(&receipt); err != nil {
		h.log.Error("Validation failed", zap.Error(err))

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			h.log.Error("Validation errors", zap.Any("errors", validationErrs))
			return c.JSON(http.StatusBadRequest, ierrors.NewValidationErrorResponse(validationErrs))
		}

		return c.JSON(http.StatusBadRequest, ierrors.NewErrorResponse(http.StatusBadRequest, err.Error()))
	}

	id, err := h.receiptProcessor.ProcessReceipt(c.Request().Context(), receipt)
	if err != nil {
		h.log.Error("Error processing receipt", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ierrors.NewErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, models.ProcessReceiptResponse{ID: id})
}

func (h *receiptHandler) GetReceiptsIdPoints(c echo.Context) error {
	h.log.Info("Getting points")

	id := c.Param("id")

	if id == "" {
		h.log.Error("Missing id parameter")
		return c.JSON(http.StatusBadRequest, ierrors.NewErrorResponse(http.StatusBadRequest, "Missing id parameter"))
	}

	points, err := h.receiptProcessor.GetPoints(c.Request().Context(), id)
	if err != nil {
		h.log.Error("Error getting points", zap.Error(err))
		if errors.Is(err, ierrors.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ierrors.NewErrorResponse(http.StatusNotFound, "No receipt found for that ID"))
		}
		return c.JSON(http.StatusInternalServerError, ierrors.NewErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, models.GetReceiptPointsResponse{Points: points})
}
