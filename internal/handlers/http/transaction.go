package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rtzgod/EWallet/internal/handlers"
	"net/http"
)

type ReceiverWalletForm struct {
	ReceiverId string  `json:"to"  binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
}

// @Summary Send Money
// @Tags Transaction
// @Description Sends money from one wallet to another
// @Accept json
// @Produce json
// @Param walletId path string true "WalletId"
// @Param transaction body ReceiverWalletForm true "Transaction"
// @Router /api/v1/wallet/{walletId}/send [post]
func (h *Handler) SendMoney(c *gin.Context) {

	var input ReceiverWalletForm
	senderId := c.Param("walletId")
	if err := c.BindJSON(&input); err != nil {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "json body incorrect")
		return
	}
	if input.Amount < 10 {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "minimal amount is 10")
		return
	}
	if senderId == input.ReceiverId {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "self transaction")
		return
	}

	senderWallet, err := h.services.Wallet.GetById(senderId)
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusNotFound, "sender wallet not found")
		return
	}
	if senderWallet.Balance < input.Amount {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "sender wallet balance not enough")
		return
	}
	if _, err := h.services.Wallet.GetById(input.ReceiverId); err != nil {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "receiver wallet not found")
		return
	}
	if err := h.services.Update(senderId, input.ReceiverId, input.Amount); err != nil {
		handlers.NewErrorResponse(c, http.StatusInternalServerError, "unable to update wallets")
		return
	}
	if err := h.services.Transaction.Create(senderId, input.ReceiverId, input.Amount); err != nil {
		handlers.NewErrorResponse(c, http.StatusInternalServerError, "unable to add transaction")
		return
	}
	c.JSON(http.StatusOK, handlers.StatusResponse{
		"ok",
	})
}

// @Summary Get History
// @Tags Transaction
// @Description returns history of transactions of wallet by id
// @Produce json
// @Param walletId path string true "WalletId"
// @Router /api/v1/wallet/{walletId}/history [get]
func (h *Handler) GetHistory(c *gin.Context) {
	id := c.Param("walletId")
	if _, err := h.services.Wallet.GetById(id); err != nil {
		handlers.NewErrorResponse(c, http.StatusNotFound, "wallet not found")
		return
	}
	transactions, err := h.services.Transaction.GetAllById(id)
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusInternalServerError, "unable to get history of transactions")
		return
	}
	c.JSON(http.StatusOK, transactions)
}
