package es

import (
	"github.com/zeromicro/go-zero/core/conf"
	"net/http"
	"time"
)

type Config struct {
	EsConf EsConfig
}
type EsConfig struct {
	Community string
}

type EsClient struct {
	*http.Client
	BaseURL string
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetEsCommunityClient() *EsClient {
	return &EsClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: c.EsConf.Community,
	}
}