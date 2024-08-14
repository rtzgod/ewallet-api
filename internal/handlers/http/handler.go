package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rtzgod/ewallet-api/internal/domain/service"

	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/rtzgod/ewallet-api/docs"
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
		wallet.POST("/", h.CreateWallet)
		wallet.GET("/:walletId", h.GetWalletById)
		wallet.POST("/:walletId/send", h.SendMoney)
		wallet.GET("/:walletId/history", h.GetHistory)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
