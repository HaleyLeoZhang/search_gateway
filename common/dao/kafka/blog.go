package kafka

import (
	"context"
	"encoding/json"
	"github.com/HaleyLeoZhang/go-component/driver/xmetric"
	"search_gateway/common/constant"
	"search_gateway/common/model/kafka"
)

// 发送博客主搜ES消息
func (d *Dao) NotifyBlogSearch(ctx context.Context, id int64, action constant.EsAction, source string) (err error) {
	err = checkNotify(id, source)
	if err != nil {
		return
	}
	var topic = constant.KAFKA_TOPIC_BLOG_SEARCH
	// 记录指标
	xmetric.MetricProducer.WithLabelValues(topic).Inc()
	// -
	msg := &kafka.BlogSearch{}
	msg.IniBaseData()
	msg.Id = id
	msg.Action = action
	msg.Source = getSourceNameForLogCenter(source)
	bs, _ := json.Marshal(msg)
	err = d.producer.SendMsgAsyncByKey(topic, msg.GetIdString(), bs)
	if err != nil {
		return
	}
	return
}
