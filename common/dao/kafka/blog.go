package kafka

import (
	"context"
	"encoding/json"
	"search_gateway/common/constant"
	"search_gateway/common/model/kafka"
)

// 发送配色关联颜色任务ES消息
func (d *Dao) NotifyBlogSearch(ctx context.Context, id int64, action constant.EsAction, source string) (err error) {
	err = checkNotify(id, source)
	if err != nil {
		return
	}
	msg := &kafka.BlogSearch{}
	msg.IniBaseData()
	msg.Id = id
	msg.Action = action
	msg.Source = getSourceNameForLogCenter(source)
	bs, _ := json.Marshal(msg)
	err = d.producer.SendMsgAsyncByKey(constant.KAFKA_TOPIC_BLOG_SEARCH, msg.GetIdString(), bs)
	if err != nil {
		return
	}
	return
}
