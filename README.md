# community
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