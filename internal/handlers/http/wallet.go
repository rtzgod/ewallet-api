package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rtzgod/EWallet/internal/handlers"
	"net/http"
)

// @Summary CreateWallet
// @Tags Wallet
// @Description Creates wallet and returns id and balance of wallet
// @Success 200 {object} handlers.WalletResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/v1/wallet/ [post]
func (h *Handler) createWallet(c *gin.Context) {
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

// @Summary GetWallet
// @Tags Wallet
// @Description Returns wallet by id
// @Param walletId path string true "WalletId"
// @Router /api/v1/wallet/{walletId} [get]
func (h *Handler) getWalletById(c *gin.Context) {
	id := c.Param("walletId")
	wallet, err := h.services.Wallet.GetById(id)
	if err != nil {
		handlers.NewErrorResponse(c, http.StatusNotFound, "unable to get wallet")
		return
	}
	c.JSON(http.StatusOK, handlers.WalletResponse{
		Id:      wallet.Id,
		Balance: wallet.Balance,
	})
}
