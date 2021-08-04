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
	"search_gateway/common/model/vo"
)

type Blog struct{}

func (*Blog) Front(c *gin.Context) {
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

	//g := &singleflight.Group{}
	//groupKey := httpSingleFightKey("blog_front", param)
	//res, err := g.Do(groupKey, func() (data interface{}, errBusiness error) {
	//	data, errBusiness = srv.CommonService.BlogFrontSearch(xGin.C, param)
	//	return
	//})
	res, err:= srv.CommonService.BlogFrontSearch(xGin.C, param)

	xGin.Response(err, res)
}
