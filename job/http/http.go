package http

import (
	"github.com/gin-gonic/gin"
	"search_gateway/job/service"
)

var srv *service.Service

func Init(e *gin.Engine, srvInjection *service.Service) *gin.Engine {
	srv = srvInjection

	return e
}
