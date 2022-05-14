package service

import (
	"encoding/json"
	"search_gateway/common/model/vo"
	"testing"
	"time"
)

func TestService_BlogFrontSearch(t *testing.T) {
	req := &vo.BlogFrontRequest{}
	keyword := "消息队列"
	req.Title = keyword
	req.Describe = keyword
	req.Category = keyword
	list, err := svr.BlogFrontSearch(ctx, req)
	if err != nil {
		t.Fatalf("Err(%+v)", err)
	}
	b, _ := json.Marshal(list)
	t.Logf("list(%v)", string(b))
}

func TestService_blogSearchEsSendAll(t *testing.T) {
	svr.blogSearchEsSendAll(ctx)
	<-time.After(3 * time.Second) // 异步发送消息
}

func TestService_blogSearchAssemble(t *testing.T) {
	var (
		id int64 = 10
	)
	gotDoc, err := svr.blogSearchAssemble(ctx, id)
	if err != nil {
		t.Fatalf("Err(%+v)", err)
	}
	b, _ := json.Marshal(gotDoc)
	t.Logf("%v", string(b))
}
