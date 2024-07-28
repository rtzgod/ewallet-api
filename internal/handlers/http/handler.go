package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rtzgod/EWallet/internal/domain/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()
	wallet := r.Group("/api/v1/wallet")
	{
		wallet.POST("/", h.createWallet)
		wallet.GET("/:walletId", h.getWalletById)
		wallet.POST("/:walletId/send", h.sendMoney)
		wallet.GET("/:walletId/history", h.getHistory)
	}
	return r
}
