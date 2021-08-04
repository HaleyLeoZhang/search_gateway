package service

import (
	"context"
	"flag"
	"os"
	"search_gateway/api/conf"
	"testing"
)

var (
	svr *Service
	ctx context.Context
)

func TestMain(m *testing.M) {
	flag.Parse()
	err := conf.Init()
	if err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
	ctx = context.Background()
	ctx = context.WithValue(ctx, "page", 1)
	os.Exit(m.Run())
}
