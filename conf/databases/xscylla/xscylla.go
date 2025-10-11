package xscylla

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	ScyllaConf ScyllaConfig
}

type ScyllaConfig struct {
	Hosts       []string // 集群节点列表
	Keyspace    string   // 键空间
	Consistency string   // 一致性级别
	Timeout     int      // 超时时间(秒)
}

var config Config

// 注意：建议通过环境变量注入配置文件路径，避免硬编码
var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"

func init() {
	// 初始化配置
	conf.MustLoad(configFile, &config)
}

// 构建集群配置
func buildClusterConfig(keyspace string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(config.ScyllaConf.Hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = parseConsistency(config.ScyllaConf.Consistency)

	timeout := time.Duration(config.ScyllaConf.Timeout) * time.Second
	cluster.Timeout = timeout
	cluster.ConnectTimeout = timeout
	cluster.DisableInitialHostLookup = true // Docker环境适用

	return cluster
}

// 解析一致性级别
func parseConsistency(consistency string) gocql.Consistency {
	switch consistency {
	case "any":
		return gocql.Any
	case "one":
		return gocql.One
	case "two":
		return gocql.Two
	case "three":
		return gocql.Three
	case "quorum":
		return gocql.Quorum
	case "all":
		return gocql.All
	default:
		return gocql.Quorum // 默认使用quorum
	}
}

func GetScyllaCommunitySession() *gocql.Session {
	cluster := buildClusterConfig(config.ScyllaConf.Keyspace)
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return session
}
