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
			name: "InternalServerError",
			mockBehavior: func(mockService *mock_service.MockWallet) {
				mockService.EXPECT().Create().Return(entity.Wallet{}, errors.New("some error"))
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
