package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rtzgod/EWallet/internal/handlers"
	"net/http"
)

func (h *Handler) createWallet(c *gin.Context) {
	wallet, err := h.services.Wallet.Create()
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusInternalServerError, "unable to create wallet")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Id":     wallet.Id,
		"Amount": wallet.Balance,
	})
}

func (h *Handler) getWalletById(c *gin.Context) {
	id := c.Param("walletId")
	wallet, err := h.services.Wallet.GetById(id)
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusNotFound, "unable to get wallet")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Id":      wallet.Id,
		"Balance": wallet.Balance,
	})
}
