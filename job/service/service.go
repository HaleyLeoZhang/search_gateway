package service

import (
	"context"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	comonnconf "search_gateway/common/conf"
	commonservice "search_gateway/common/service"
	"search_gateway/job/conf"
)

type Service struct {
	cfg           *conf.Config
	commonService *commonservice.Service
	// 全局上下文
	ctx       context.Context // 带取消，等关闭进程的时候，自动触发关闭
	ctxCancel context.CancelFunc
}

// New create service instance and return.
func New(cfg *conf.Config) *Service {
	s := &Service{}
	s.cfg = cfg
	// 初始化: 配置
	cfgCommon := &comonnconf.Config{}
	cfgCommon.DB = cfg.DB
	cfgCommon.Redis = cfg.Redis
	cfgCommon.Es = cfg.Es
	cfgCommon.Kafka = cfg.Kafka
	s.commonService = commonservice.New(cfgCommon)
	// 初始化: 全局上下文
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	// 初始
	return s
}

func (s *Service) Start() {
	s.initConsumers()
}

// Close close the resource.
func (s *Service) Close() {
	// 通知全局上文关闭
	s.ctxCancel()
	// 各种消费者
	// - 暂无
	// 各种数据库
	// - 平滑关闭，建议数据库相关的关闭放到最后
	s.commonService.Close()
	xlog.Info("Close.Service.Done")
}
