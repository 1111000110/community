package xscylla

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	ScyllaConf ScyllaConfig `yaml:"ScyllaConf"`
}

type ScyllaConfig struct {
	Hosts       []string `yaml:"Hosts"`       // 集群节点列表
	Keyspace    string   `yaml:"Keyspace"`    // 键空间
	Consistency string   `yaml:"Consistency"` // 一致性级别
	Timeout     int      `yaml:"Timeout"`     // 超时时间(秒)
}

var config Config

// 注意：建议通过环境变量注入配置文件路径，避免硬编码
var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"

func init() {
	// 初始化配置
	conf.MustLoad(configFile, &config)

	// 提前检查并创建Keyspace（只在初始化时执行一次）
	if err := ensureKeyspaceExists(); err != nil {
		fmt.Printf("Warning: failed to ensure keyspace exists: %v\n", err)
	}
}

// 确保Keyspace存在，不存在则创建
func ensureKeyspaceExists() error {
	// 先尝试连接目标Keyspace
	cluster := buildClusterConfig(config.ScyllaConf.Keyspace)
	session, err := cluster.CreateSession()
	if err == nil {
		session.Close()
		return nil
	}

	// 检查是否是Keyspace不存在的错误
	if isKeyspaceNotFoundError(err) {
		fmt.Printf("Keyspace %s not found, creating it...\n", config.ScyllaConf.Keyspace)

		// 连接系统Keyspace创建目标Keyspace
		systemCluster := buildClusterConfig("system")
		systemSession, err := systemCluster.CreateSession()
		if err != nil {
			return fmt.Errorf("failed to connect to system keyspace: %v", err)
		}
		defer systemSession.Close()

		// 创建Keyspace
		createStmt := fmt.Sprintf(`
			CREATE KEYSPACE IF NOT EXISTS %s 
			WITH REPLICATION = {
				'class': 'SimpleStrategy', 
				'replication_factor': 1
			}
		`, config.ScyllaConf.Keyspace)

		return systemSession.Query(createStmt).Exec()
	}

	return fmt.Errorf("unexpected error checking keyspace: %v", err)
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

// 判断错误是否为Keyspace不存在
func isKeyspaceNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// 匹配Scylla返回的Keyspace不存在错误信息
	return gocql.ErrNoKeyspace.Error() == err.Error() ||
		fmt.Sprintf("Keyspace %s does not exist", config.ScyllaConf.Keyspace) == err.Error()
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
