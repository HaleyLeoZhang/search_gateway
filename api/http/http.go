package http

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"search_gateway/api/service"
)

var srv *service.Service

func Init(e *gin.Engine, srvInjection *service.Service) *gin.Engine {
	srv = srvInjection

	{
		blog := &Blog{}
		i := e.Group("blog/")
		i.GET("front/", blog.Front) // 前台搜素
	}

	return e
}

func httpSingleFightKey(route string, param interface{}) (str string) {
	bytes, _ := json.Marshal(param)
	md5Raw := md5.Sum(bytes)
	md5String := fmt.Sprintf("%x", md5Raw)
	str = fmt.Sprintf("%v:md5:%v:v1", route, md5String)
	return
}
