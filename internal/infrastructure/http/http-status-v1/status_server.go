package http_status_v1

import (
	"github.com/gin-gonic/gin"

	"github.com/TestardR/geo-tracking/config"
)

type getStatusHttpHandler interface {
	GetStatus(ctx *gin.Context)
}

type StatusServer struct {
	router           *gin.Engine
	port             string
	getStatusHandler getStatusHttpHandler
}

func NewStatusHttpServer(
	cfg *config.Config,
	handler getStatusHttpHandler,
) StatusServer {
	gin.SetMode(cfg.Env)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	r := router.Group("/v1")

	r.GET("/status/:id", handler.GetStatus)

	return StatusServer{
		router: router,
		port:   cfg.HttpPort,
	}
}
