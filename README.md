# community
## 项目部署
### 10.27部署复盘
初期本人（张璇）直接根据go-zero文档（`https://go-zero.dev/docs/tutorials/ops/prepare` ）在本机（mac）和云服务器（2h4glinux服务器）部署都没有成功，原因如下：
1. 初期希望代码在后期用户变多的时候直接能复用，而不需要改库改代码，所使用的组件过多（微服务，k8s，redis，kafka，mysql，scyllaDb，elasticSearch，gitlab等），但为减少成本使用2h4g服务器，未考虑机器性能，实际部署仅gitlab就导致服务器崩溃。
2. mac因为架构的问题，很多docker容器（elasticsearch，jenkins等）无法兼容。
后续考虑到用户量可能不多甚至没有用户，现有框架纯粹是过度设计，没有实际必要，如消息存储mysql在千万数据量下仍可使用，但我直接使用scyllaDb做分布式，导致大量问题出现且无法部署使用。
基于上述原因，采取以下措施，尽量保证后期扩展： 
- 使用Github，GitHub Actions，GitHub Packages前期替代Gitlab，Jenkins，Harbor，无需部署直接使用线上免费的软件，等流量到了收费的程度再考虑自建服务。
- 为节约成本，nginx，项目服务，以及使用的一切第三方软件（mysql，kafka，redis）等全部使用docker部署，放弃k8s，使用docker-compose部署。
- 停止使用elasticSearch和scyllaDb，后续如果必要再切换回来，项目保证解耦合，防止耦合过重导致数据库无法切换。
## 开发规范
### 目录结构
```
service/服务名
- client：服务内部和外部都直接调用的变量和函数，不允许引用任何其他非第三方包，也不允许调用任何rpc接口，仅做全局声明使用。
- api：web和客户端直接调用的接口层，一般用来鉴权和服务调用，其他所有文件不能引用api内的任何文件。
- rpc：处理业务逻辑的服务层，可调用任何rpc，也可被任何rpc调用。
- model：数据交互层，直接与数据库进行交互。
```
### api文档书写
api服务及内部的所有接口不允许其他任何服务直接调用
1. .api文件
示例：https://go-zero.dev/docs/reference
所有命名及格式问题参照：community/service/message
```
type (
    // 大驼峰定义请求体，结构为：服务名+行为+Req/Resp
	MessageUpdateByIdReq {
		MessageIds string   `json:"message_ids"`// 大驼峰定义名称，小驼峰定义json
	}
)
@server (
	group: message // 全小写命名
)
// 大驼峰命名，结构为：服务名
service Message {
    // 大驼峰命名，结构为：服务名+行为
  	@handler MessageUpdateById
  	// 下划线命名路径/服务/行为 
	post /message/update/id (MessageUpdateByIdReq) returns (MessageUpdateByIdResp)
}
```
2. config.go文件
```
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	MessageConf       zrpc.RpcClientConf // 所有调用的rpc命名直接为{服务名{大驼峰}+Conf}
	WebSocketPushConf zrpc.RpcClientConf 
}
```
3. svc目录
- 目录结构
```
/svc
- setvicecontext.go 上下文定义
- rpcclient.go 所有的rpc定义
- modelclient.go 所有的中间件及数据库定义
```
- servicecontext.go
```
// servicecontext.go
type ServiceContext struct {
	Config    config.Config
	RpcClient *RpcClient // rpc客户端，固定格式命名
	ModelClient *ModelClient // 数据库客户端，固定格式命名
}
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		RpcClient: NewRpcClient(c),
		ModelClient: DefaultModelClient(),
	}
}
```
- rpcclient.go
```
type RpcClient struct {
	Message       message.Message // 驼峰命令，格式服务名，包名也必须为服务名小写，如果正确格式构建rpc，一定得到的是不重复的小写名
	WebSocketPush websocketpush.WebsocketPush // 驼峰命令，格式服务名
}

func NewRpcClient(c config.Config) *RpcClient {
	return &RpcClient{
		Message:       message.NewMessage(zrpc.MustNewClient(c.MessageConf)),
		WebSocketPush: websocketpush.NewWebsocketPush(zrpc.MustNewClient(c.WebSocketPushConf)),
	}
}
```
- modelclient.go
```
type ModelClient struct {
    // 对于唯一端，比如scylla，mysql等，通过表名来区分的直接写数据库名称
    // 对于多端，即需要参数来区分的，写数据库+表名
	Mysql user.UserModel
	MongoPost model.PostModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		Mysql: user.NewUserModel(xmysql.GetMysqlCommunityClient()),
		MongoPost: model.NewCommunityModel("Post"),
	}
}

```
## 附件
### 附件一：命名方式
群聊创建
1. 下划线：chat_group_create
2. 小驼峰：chatGroupCreate
3. 大驼峰：ChatGroupCreate