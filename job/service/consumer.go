package service

import (
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xkafka"
	"search_gateway/common/constant"
	"time"
)

// 初始化消费任务
func (s *Service) initConsumers() {
	if s.cfg.Kafka == nil {
		fmt.Println("当前没有kafka配置文件")
		return
	}
	// 业务类型: 博客前台搜索
	_ = xkafka.StartKafkaConsumer(s.ctx, xkafka.ConsumerOption{
		Conf:        s.cfg.Kafka,                                // Kafka 配置
		Topic:       []string{constant.KAFKA_TOPIC_BLOG_SEARCH}, // 消费Topic列表
		Group:       constant.KAFKA_GROUP_BLOG_SEARCH,           // Consumer group name
		Batch:       10,                                         // 每次拉取的消息数
		Procs:       2,                                          // 并发处理消息的数量
		PollTimeout: 3 * time.Second,
		Handler:     s.commonService.KafkaBlogSearchEs, // 处理单条消息的函数
		Mode:        xkafka.ModeBatch,                  // 消费模式
	})
}
