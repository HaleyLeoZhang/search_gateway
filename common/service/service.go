package service

import (
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"search_gateway/common/conf"
	"search_gateway/common/dao/cache"
	"search_gateway/common/dao/es"
	"search_gateway/common/dao/kafka"
	"search_gateway/common/dao/yth_blog_avatar"
)

type Service struct {
	BlogDao  *yth_blog_avatar.Dao
	EsDao    *es.Dao
	CacheDao *cache.Dao
	producer *kafka.Dao
}

func New(cfg *conf.Config) *Service {
	s := &Service{}
	if cfg.Redis != nil {
		s.CacheDao = cache.New(cfg.Redis)
	}
	if cfg.Es != nil {
		s.EsDao = es.New(cfg.Es)
	}
	if cfg.DB != nil {
		s.BlogDao = yth_blog_avatar.New(cfg.DB)
	}
	if cfg.Kafka != nil {
		s.producer = kafka.New(cfg.Kafka)
	}
	return s
}

// Close close the resource.
func (s *Service) Close() {
	// 各种MQ的生产、消费者
	if s.producer != nil {
		s.producer.Close()
	}
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
