# Qelog
是一款小巧低成本轻运维的日志服务。其诞生的目的就是为了解决中小团队多服务群日志和报警问题。

已经稳定运行2年+，滚动写入超过100+TB数据。

[后台示例地址](https://qelogdemo.bbdshow.top/admin)  账号：admin 密码：111111

> 创建资源时请文明用于，谢谢合作！

### 日志系统特性:

#### Client
- **qelog** 是项目的Go版本Client。 Wrap [Uber-zap](https://github.com/uber-go/zap) ，扩展功能，对WriteSyncer进行本地和远程的实现。
- 本地扩展：日志的切割，压缩，保留期，动态切换等级等。
- 远程扩展：缓冲打包压缩传输日志数据，异常备份重试保证不丢失，支持GRPC(默认) HTTP协议，IO控制保证带宽占用可控，内存占用低等特点。

#### Log receiver server
- Receiver 进程可横向扩展，保证高可用，高性能。
- 通过对Receiver配置多存储实例，提高集群存储容量和写入性能。
- 报警模块，对每条日志进行报警规则检测。命中后可以根据规则和不同的报警方式送达。目前支持 DingTalk | Telegram
- 实现数据分片存储规则，支持自动管理容量，监控预警。存储扩展单独实例，不因中间件而产生瓶颈。
- 日志统计，等级分布，趋势报表等。

#### Log manager server
- 基于 [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin) 修改，支持丰富的查询维度，比如 等级、关键字、TraceId、ClientIP、多级条件等。因成本和性能问题**暂不支持全文索引**
- 友好的操作交互，简易的配置，高效的内容展示，几乎可做到**开箱注册即用**。
- 集群容量统计查询，手动干预等功能。

#### 性能测试工具

- /tools 包含benchmark与内存占用测试。
- 不同的运行环境和配置有不同的表现，可合理综合分析，初步检测此服务是否满足场景需求。

PS: 生产环境中，该项目经过几轮死循环“攻击”。。。

#### 设计简图

![设计简图](https://qnoss.bbdshow.top/notes/qelog.png)

### 使用建议

#### Client端导入项目

> go get github.com/bbdshow/qelog/qezap

[Client示例](../qezap/example/main.go)


#### 服务快速部署

##### 容器部署

```shell
git clone https://github.com/bbdshow/qelog.git
cd qelog
# custom you config
vim config.docker.toml
# build docker image
make image
# docker-compose start container
docker-compose up -d
```

##### 主机部署
```shell
git clone https://github.com/bbdshow/qelog.git
cd qelog

# go build, output ./bin
make
cd ./bin
vim configs/config.toml
cd ./admin
# admin suggest single server, because have background task
nohup ./qelog_admin -f ../configs/config.toml >> nohup.out 2>&1  &
nohup ./qelog_receiver -f ../configs/config.toml >> nohup.out 2>&1  &

```

感谢支持，如果对您有用，希望**Star**以表支持,有问题请提 Issues，持续更新并解决问题。

