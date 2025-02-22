package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"ticket-processor/internal/ierrors"
	"ticket-processor/internal/models"
)

type MockReceiptProcessor struct {
	mock.Mock
}

func (m *MockReceiptProcessor) ProcessReceipt(ctx context.Context, receipt models.Receipt) (string, error) {
	args := m.Called(ctx, receipt)
	return args.String(0), args.Error(1)
}

func (m *MockReceiptProcessor) GetPoints(ctx context.Context, id string) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func TestReceiptHandler_PostReceiptsProcess_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(`{
		"retailer":"Retailer",
		"purchaseDate":"2023-10-10",
		"purchaseTime":"10:10",
		"items":[{"shortDescription":"Item","price":"1.00"}],
		"total":"1.00"
	}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockProcessor := new(MockReceiptProcessor)
	mockProcessor.On("ProcessReceipt", mock.Anything, mock.Anything).Return("123", nil)

	handler := NewReceiptHandler(zap.NewNop(), mockProcessor)

	err := handler.PostReceiptsProcess(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"id":"123"}`, rec.Body.String())
}

func TestReceiptHandler_PostReceiptsProcess_InvalidJSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(`invalid json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewReceiptHandler(zap.NewNop(), new(MockReceiptProcessor))

	err := handler.PostReceiptsProcess(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"statusText":"Bad Request","message":"Invalid JSON format"}`, rec.Body.String())
}

func TestReceiptHandler_GetReceiptsIdPoints_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/receipts/123/points", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	mockProcessor := new(MockReceiptProcessor)
	mockProcessor.On("GetPoints", mock.Anything, "123").Return(100, nil)

	handler := NewReceiptHandler(zap.NewNop(), mockProcessor)

	err := handler.GetReceiptsIdPoints(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"points":100}`, rec.Body.String())
}

func TestReceiptHandler_GetReceiptsIdPoints_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/receipts/123/points", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	mockProcessor := new(MockReceiptProcessor)
	mockProcessor.On("GetPoints", mock.Anything, "123").Return(0, ierrors.ErrNotFound)

	handler := NewReceiptHandler(zap.NewNop(), mockProcessor)

	err := handler.GetReceiptsIdPoints(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"statusText":"Not Found","message":"No receipt found for that ID"}`, rec.Body.String())
}

func TestReceiptHandler_GetReceiptsIdPoints_MissingID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/receipts//points", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewReceiptHandler(zap.NewNop(), new(MockReceiptProcessor))

	err := handler.GetReceiptsIdPoints(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"statusText":"Bad Request","message":"Missing id parameter"}`, rec.Body.String())
}

func TestReceiptHandler_PostReceiptsProcess_MissingFields(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(`{"retailer":""}`)) // Missing required fields
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockProcessor := new(MockReceiptProcessor)
	handler := NewReceiptHandler(zap.NewNop(), mockProcessor)

	err := handler.PostReceiptsProcess(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"statusText":"Bad Request","message":"Invalid JSON format"}`, rec.Body.String())

}
