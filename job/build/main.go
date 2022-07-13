package main

import (
	"flag"
	"github.com/HaleyLeoZhang/go-component/driver/bootstrap"
	"github.com/HaleyLeoZhang/go-component/driver/httpserver"
	"github.com/HaleyLeoZhang/go-component/driver/xgin"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"search_gateway/job/conf"
	"search_gateway/job/http"
	"search_gateway/job/service"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	oneService := service.New(conf.Conf)
	xlog.Init(conf.Conf.Log)

	gin := xgin.New(conf.Conf.Gin)

	app := bootstrap.New()
	app.Start(func() { // 此部分代码，请勿阻塞进程
		oneService.Start()
		go httpserver.Run(conf.Conf.HttpServer, http.Init(gin, oneService)) // 已配置 recovery 不用处理 panic
		return
	}).Stop(func() {
		oneService.Close()
	})

}
