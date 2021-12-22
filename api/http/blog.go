package http

// ----------------------------------------------------------------------
// 漫画控制器
// ----------------------------------------------------------------------
// Link  : http://www.hlzbxlog.top/
// GITHUB: https://github.com/HaleyLeoZhang
// ----------------------------------------------------------------------
import (
	"github.com/HaleyLeoZhang/go-component/driver/xgin"
	"github.com/gin-gonic/gin"
	"github.com/golang/groupcache/singleflight"
	"search_gateway/common/model/vo"
)

type Blog struct {
	singleFlightFront singleflight.Group // 接口级缓存 幂等请求，防止击穿 说明文档 https://segmentfault.com/a/1190000018464029
}

func (b *Blog) Front(c *gin.Context) {
	xGin := xgin.NewGin(c)
	var (
		err   error
		param = &vo.BlogFrontRequest{}
	)

	err = c.Bind(param)
	if err != nil {
		err = &xgin.BusinessError{Code: xgin.HTTP_RESPONSE_CODE_PARAM_INVALID, Message: "Param is invalid"}
		xGin.Response(err, nil)
		return
	}

	// 幂等请求，防止击穿 说明文档 https://segmentfault.com/a/1190000018464029
	groupKey := httpSingleFightKey("blog_front", param)
	res, err := b.singleFlightFront.Do(groupKey, func() (data interface{}, errBusiness error) {
		data, errBusiness = srv.CommonService.BlogFrontSearch(xGin.C, param)
		return
	})
	xGin.Response(err, res)
}
