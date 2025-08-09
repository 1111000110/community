package mongo

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
)

type Config struct {
	MongoConf MongoConfig
}
type MongoConfig struct {
	Community string
	Db        string
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetMongoCommunityClient(collection string) *mon.Model { // mongo的可生产文件会依赖此函数，尽量不要修改位置
	model, err := mon.NewModel(c.MongoConf.Community, c.MongoConf.Db, collection)
	logx.Must(err)
	return model
}
