package service

import (
	"context"
	"os"
	"search_gateway/common/conf"
	"testing"
)

var (
	svr *Service
	ctx context.Context
)

func TestMain(m *testing.M) {
	err := conf.Init() // 注：可以参照 common/conf/debug.example.yaml 新建一个 common/conf/debug.yaml
	if err != nil {
		panic(err)
	}
	cfg := conf.Conf
	svr = New(cfg)
	ctx = context.Background()
	os.Exit(m.Run())
}
