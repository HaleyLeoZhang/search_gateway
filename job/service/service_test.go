package service

import (
	"context"
	"flag"
	"os"
	"search_gateway/job/conf"
	"testing"
)

var (
	svr *Service
	ctx context.Context
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../build/app.yaml")
	err := conf.Init()
	if err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
	ctx = context.Background()
	ctx = context.WithValue(ctx, "page", 1)
	ctx = context.WithValue(ctx, "page_id", 1311)
	os.Exit(m.Run())
}
