package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rtzgod/ewallet-api/internal/handlers"
	"net/http"
)

// @Summary Create Wallet
// @Tags Wallet
// @Description Creates wallet and returns id and balance of wallet
// @Success 200 {object} handlers.WalletResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/v1/wallet/ [post]
func (h *Handler) CreateWallet(c *gin.Context) {
	wallet, err := h.services.Wallet.Create()
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusInternalServerError, "unable to create wallet")
		return
	}
	c.JSON(http.StatusOK, handlers.WalletResponse{
		Id:      wallet.Id,
		Balance: wallet.Balance,
	})
}

// @Summary Get Wallet
// @Tags Wallet
// @Description Returns wallet by id
// @Param walletId path string true "WalletId"
// @Router /api/v1/wallet/{walletId} [get]
func (h *Handler) GetWalletById(c *gin.Context) {
	id := c.Param("walletId")
	wallet, err := h.services.Wallet.GetById(id)
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusNotFound, "wallet not found")
		return
	}
	c.JSON(http.StatusOK, handlers.WalletResponse{
		Id:      wallet.Id,
		Balance: wallet.Balance,
	})
}
