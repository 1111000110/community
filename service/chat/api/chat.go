package main

import (
	"community.com/chat/api/internal/config"
	"community.com/chat/api/internal/handler"
	"community.com/chat/api/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/chat.yaml", "the config file")

func main() {
	flag.Parse()
	logx.Disable() // 禁用 go-zero 的默认日志
	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	go ctx.Hub.Run()       //启动聊天中心
	go ctx.Receive.Start() //启动消息队列监听
	defer ctx.Receive.Stop()
	server.Start()
}
