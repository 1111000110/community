package mysql

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	MysqlConf MysqlConfig
}
type MysqlConfig struct {
	Community string
}

var configFile = flag.String("f", "conf.yaml", "the config file")
var c Config

func init() {
	flag.Parse()
	conf.MustLoad(*configFile, &c)
}

func GetMysqlCommunityClient() sqlx.SqlConn {
	return sqlx.NewMysql(c.MysqlConf.Community)
}
