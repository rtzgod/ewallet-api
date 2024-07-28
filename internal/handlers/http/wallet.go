package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createWallet(c *gin.Context) {
	wallet, err := h.services.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Id":     wallet.Id,
		"Amount": wallet.Amount,
	})
}

func (h *Handler) getWalletById(c *gin.Context) {
	id := c.Param("walletId")
	wallet, err := h.services.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"wallet not found, error:": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Id":      wallet.Id,
		"Balance": wallet.Amount,
	})
}
