package kafka

import (
	"fmt"
	"github.com/HaleyLeoZhang/go-component/driver/xkafka"
	"github.com/pkg/errors"
)

type Dao struct {
	producer *xkafka.Producer
}

func New(cfg *xkafka.Config) *Dao {
	d := &Dao{}
	d.producer = xkafka.NewProducer(cfg)
	return d
}

func (d *Dao) Close() {
	_ = d.producer.Close()
}

// 以下是一些校验逻辑

func checkNotify(id int64, source string) (err error) {
	if id == 0 {
		err = errors.WithStack(fmt.Errorf("ID必须大于0"))
		return
	}
	if source == "" {
		err = errors.WithStack(fmt.Errorf("source 不能为空"))
		return
	}
	return
}

// 获取日志中心专属的Source开头
func getSourceNameForLogCenter(source string) string {
	return fmt.Sprintf("searh_gateway_%v", source)
}
