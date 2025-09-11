package main

import (
	"community/service/websocket/api/internal/config"
	"community/service/websocket/api/internal/handler"
	"community/service/websocket/api/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/websocket.yaml", "the config file")

// todo 可能不安全 记得处理
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	serviceGroup := service.NewServiceGroup() // 服务组
	defer serviceGroup.Stop()

	server := rest.MustNewServer(c.RestConf)
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	serviceGroup.Add(ctx.MessageHub)       // 消息中心服务
	for _, mq := range ctx.KafkaConsumer { // kafka消费组
		serviceGroup.Add(mq)
	}
	serviceGroup.Add(server) // http服务

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	serviceGroup.Start()
}
