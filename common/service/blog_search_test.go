package service

import (
	"encoding/json"
	"search_gateway/common/model/vo"
	"testing"
)

func TestService_BlogFrontFlushAll(t *testing.T) {
	svr.BlogFrontFlushAll()
}

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

func TestService_BlogFrontIni(t *testing.T) {
	err := svr.BlogFrontIni()
	if err != nil {
		t.Fatalf("Err(%+v)", err)
	}
}
