package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	WebSocket struct {
		WriteWait      int64
		PongWait       int64
		PingPeriod     int64
		MaxMessageSize int64
		BufSize        int64
	}
}
