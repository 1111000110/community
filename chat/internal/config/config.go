package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Receive kq.KqConf
	Pusher  struct {
		Brokers []string
		Topic   string
	}
	WebSocket struct {
		WriteWait      int64
		PongWait       int64
		PingPeriod     int64
		MaxMessageSize int64
		BufSize        int64
	}
}
