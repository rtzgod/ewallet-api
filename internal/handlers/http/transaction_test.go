package http

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rtzgod/ewallet-api/internal/domain/entity"
	"github.com/rtzgod/ewallet-api/internal/domain/service"
	mock_service "github.com/rtzgod/ewallet-api/internal/domain/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_sendMoney(t *testing.T) {
	type MockBehavior func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm)

	testTable := []struct {
		name                string
		walletId            string
		inputBody           string
		inputReceiver       ReceiverWalletForm
		mockBehavior        MockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			walletId:  "1",
			inputBody: `{"to": "2", "amount": 10}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "2",
				Amount:     10,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 100}, nil)
				mockWalletService.EXPECT().GetById(receiver.ReceiverId).Return(entity.Wallet{Id: "2", Balance: 100}, nil)
				mockWalletService.EXPECT().Update(walletId, receiver.ReceiverId, receiver.Amount).Return(nil)
				mockTransactionService.EXPECT().Create(walletId, receiver.ReceiverId, receiver.Amount).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:      "non-existent sender wallet",
			walletId:  "this wallet doesn't exist",
			inputBody: `{"to": "2", "amount": 10}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "2",
				Amount:     10,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{}, errors.New("wallet not found"))
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"sender wallet not found"}`,
		},
		{
			name:      "non-existent receiver wallet",
			walletId:  "1",
			inputBody: `{"to": "this wallet doesn't exist", "amount": 10}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "this wallet doesn't exist",
				Amount:     10,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 100}, nil)
				mockWalletService.EXPECT().GetById(receiver.ReceiverId).Return(entity.Wallet{}, errors.New("wallet not found"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"receiver wallet not found"}`,
		},
		{
			name:      "incorrect json body",
			walletId:  "1",
			inputBody: `{"wrong json key": "2", "amount": 10}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "2",
				Amount:     10,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"json body incorrect"}`,
		},
		{
			name:      "same sender and receiver wallet ID's",
			walletId:  "1",
			inputBody: `{"to": "1", "amount": 10}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "1",
				Amount:     10,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"self transaction"}`,
		},
		{
			name:      "amount of sending money less than 10",
			walletId:  "1",
			inputBody: `{"to": "2", "amount": -100}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "1",
				Amount:     -100,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"minimal amount is 10"}`,
		},
		{
			name:      "not enough money on sender balance",
			walletId:  "1",
			inputBody: `{"to": "2", "amount": 110}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "1",
				Amount:     110,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 0}, nil)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"sender wallet balance not enough"}`,
		},
		{
			name:      "Internal DB Error while updating wallets info",
			walletId:  "1",
			inputBody: `{"to": "2", "amount": 10}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "2",
				Amount:     10,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 100}, nil)
				mockWalletService.EXPECT().GetById(receiver.ReceiverId).Return(entity.Wallet{Id: "2", Balance: 100}, nil)
				mockWalletService.EXPECT().Update(walletId, receiver.ReceiverId, receiver.Amount).Return(errors.New("some db error while updating wallets"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"unable to update wallets"}`,
		},
		{
			name:      "Internal DB Error while adding transaction info to transactions table",
			walletId:  "1",
			inputBody: `{"to": "2", "amount": 10}`,
			inputReceiver: ReceiverWalletForm{
				ReceiverId: "2",
				Amount:     10,
			},
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string, receiver ReceiverWalletForm) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 100}, nil)
				mockWalletService.EXPECT().GetById(receiver.ReceiverId).Return(entity.Wallet{Id: "2", Balance: 100}, nil)
				mockWalletService.EXPECT().Update(walletId, receiver.ReceiverId, receiver.Amount).Return(nil)
				mockTransactionService.EXPECT().Create(walletId, receiver.ReceiverId, receiver.Amount).Return(errors.New("some db error while adding transaction"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"unable to add transaction"}`,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransaction(ctrl)
	mockWalletService := mock_service.NewMockWallet(ctrl)
	services := &service.Service{Wallet: mockWalletService, Transaction: mockTransactionService}
	handler := NewHandler(services)

	r := gin.Default()
	r.POST("/api/v1/wallet/:walletId/send", handler.SendMoney)

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mockTransactionService, mockWalletService, tc.walletId, tc.inputReceiver)

			req, _ := http.NewRequest("POST", "/api/v1/wallet/"+tc.walletId+"/send", bytes.NewBufferString(tc.inputBody))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getHistory(t *testing.T) {
	type MockBehavior func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string)

	testTable := []struct {
		name                string
		walletId            string
		mockBehavior        MockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:     "OK",
			walletId: "1",
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 100}, nil)
				mockTransactionService.EXPECT().GetAllById(walletId).Return([]entity.Transaction{{SenderId: walletId, ReceiverId: "2", Amount: 10}}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `[{"time":"0001-01-01T00:00:00Z","from":"1","to":"2","amount":10}]`,
		},
		{
			name:     "Wallet doesn't exist",
			walletId: "this wallet doesn't exist",
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{}, errors.New("wallet not found"))
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"wallet not found"}`,
		},
		{
			name:     "Internal DB Error while fetching history",
			walletId: "1",
			mockBehavior: func(mockTransactionService *mock_service.MockTransaction, mockWalletService *mock_service.MockWallet, walletId string) {
				mockWalletService.EXPECT().GetById(walletId).Return(entity.Wallet{Id: "1", Balance: 100}, nil)
				mockTransactionService.EXPECT().GetAllById(walletId).Return([]entity.Transaction{}, errors.New("some db error while fetching history"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"unable to get history of transactions"}`,
		}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mock_service.NewMockTransaction(ctrl)
	mockWalletService := mock_service.NewMockWallet(ctrl)

	services := &service.Service{Wallet: mockWalletService, Transaction: mockTransactionService}
	handler := NewHandler(services)

	r := gin.Default()
	r.GET("/api/v1/wallet/:walletId/history", handler.GetHistory)

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mockTransactionService, mockWalletService, tc.walletId)

			req, _ := http.NewRequest("GET", "/api/v1/wallet/"+tc.walletId+"/history", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}
