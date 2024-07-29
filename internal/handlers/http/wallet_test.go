package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/domain/service"
	mock_service "github.com/rtzgod/EWallet/internal/domain/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_createWallet(t *testing.T) {
	type mockBehavior func(s *mock_service.MockWallet)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(mockService *mock_service.MockWallet) {
				mockService.EXPECT().Create().Return(entity.Wallet{Id: "1", Balance: 100}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":"1","balance":100}`,
		},
		{
			name: "DB Internal Error, can't create a wallet",
			mockBehavior: func(mockService *mock_service.MockWallet) {
				mockService.EXPECT().Create().Return(entity.Wallet{}, errors.New("can't create a wallet"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"unable to create wallet"}`,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletService := mock_service.NewMockWallet(ctrl)

	services := &service.Service{Wallet: mockWalletService}
	handler := NewHandler(services)

	r := gin.Default()
	r.POST("/api/v1/wallet/", handler.createWallet)

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockWalletService)

			req, _ := http.NewRequest("POST", "/api/v1/wallet/", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getWalletById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockWallet, walletId string)

	testTable := []struct {
		name                string
		walletId            string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:     "OK",
			walletId: "1",
			mockBehavior: func(mockService *mock_service.MockWallet, walletId string) {
				mockService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 100}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":"1","balance":100}`,
		},
		{
			name:     "non-existent wallet",
			walletId: "this wallet doesn't exist",
			mockBehavior: func(mockService *mock_service.MockWallet, walletId string) {
				mockService.EXPECT().GetById(walletId).Return(entity.Wallet{}, errors.New("wallet not found"))
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"wallet not found"}`,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletService := mock_service.NewMockWallet(ctrl)
	services := &service.Service{Wallet: mockWalletService}
	handler := NewHandler(services)

	r := gin.Default()
	r.GET("/api/v1/wallet/:walletId", handler.getWalletById)

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockWalletService, tt.walletId)

			req, _ := http.NewRequest("GET", "/api/v1/wallet/"+tt.walletId, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}
