package service

import (
	"github.com/HaleyLeoZhang/go-component/driver/xelastic"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"search_gateway/common/conf"
	"search_gateway/common/dao/cache"
	"search_gateway/common/dao/yth_blog_avatar"

	v7 "github.com/olivere/elastic/v7"
)

type Service struct {
	BlogDao  *yth_blog_avatar.Dao
	EsDao    *v7.Client
	CacheDao *cache.Dao
}

func New(cfg *conf.Config) *Service {
	s := &Service{}
	if cfg.Redis != nil {
		s.CacheDao = cache.New(cfg.Redis)
	}
	if cfg.Es != nil {
		s.EsDao, _ = xelastic.NewV7(cfg.Es)
	}
	if cfg.DB != nil {
		s.BlogDao = yth_blog_avatar.New(cfg.DB)
	}
	return s
}

// Close close the resource.
func (s *Service) Close() {
	// 各种消费者
	// - 暂无
	// 各种数据库
	// - 平滑关闭，建议数据库相关的关闭放到最后
	if s.CacheDao != nil {
		s.CacheDao.Close()
	}
	if s.BlogDao != nil {
		s.BlogDao.Close()
	}
	xlog.Info("Close.commonService.Done")
}
