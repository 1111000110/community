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