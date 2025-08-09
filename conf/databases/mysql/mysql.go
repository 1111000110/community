package mysql

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	MysqlConf MysqlConfig
}
type MysqlConfig struct {
	Community string
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetMysqlCommunityClient() sqlx.SqlConn {
	return sqlx.NewMysql(c.MysqlConf.Community)
}
