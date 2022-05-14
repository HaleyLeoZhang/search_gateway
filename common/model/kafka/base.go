package kafka

import (
	"fmt"
	"github.com/google/uuid"
	"search_gateway/common/constant"
	"time"
)

// 公用消息体

type kafkaMessageBase struct {
	Id     int64             `json:"id"`
	Action constant.EsAction `json:"action"`
	// 下面的是方便查数据用的
	UniqueId string `json:"unique_id"`
	Time     string `json:"time"`
	Source   string `json:"source"`
}

func (baseMsg *kafkaMessageBase) GetIdString() string {
	return fmt.Sprintf("%v", baseMsg.Id)
}

// 初始化 补充信息

func (baseMsg *kafkaMessageBase) IniBaseData() {
	baseMsg.UniqueId = uuid.New().String()
	baseMsg.Time = time.Now().Format(constant.BUSINESS_SHOW_TIME_TPL)
}
