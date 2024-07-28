package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/handlers"
	"net/http"
)

func (h *Handler) sendMoney(c *gin.Context) {
	var input entity.Transaction
	senderId := c.Param("walletId")
	if err := c.BindJSON(&input); err != nil {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "json body incorrect")
		return
	}

	if senderId == input.ReceiverId {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "self transaction")
		return
	}

	senderWallet, err := h.services.Wallet.GetById(senderId)
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusNotFound, "sender Wallet not found")
		return
	}
	if _, err := h.services.Wallet.GetById(input.ReceiverId); err != nil {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "receiver Wallet not found")
		return
	}
	if senderWallet.Balance < input.Amount {
		handlers.NewErrorResponse(c, http.StatusBadRequest, "sender wallet balance not enough")
		return
	}
	if err := h.services.Update(senderId, input.ReceiverId, input.Amount); err != nil {
		handlers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.Transaction.Create(senderId, input.ReceiverId, input.Amount); err != nil {
		handlers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, handlers.StatusResponse{
		"ok",
	})
}

func (h *Handler) getHistory(c *gin.Context) {

}
