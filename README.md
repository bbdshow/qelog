# Qelog

Qelog是一款小巧且功能干练实用的日志系统(集成报警)。采用Push模式，拥有优秀写入速度、可存储容量高，占用资源少，部署维护成本低等特性。其诞生的目的就是为了解决在中小团队或个人Go项目中，避免为日志系统投入过多的使用和运维成本(ELK)。

采用简单的设计做日志系统最本质的事情。已系统高效稳定运行、低使用成本及服务运维成本为目标。

PS:目前已经过上千TB数据写入稳定性验证，不足之处目前Client端只支持Uber-zap包, 支持写入协议(api/receiverpb)，可自行定制Client

### 日志系统特性:

1. Client端采用 Uber-zap 的Format，高效格式化日志，友好迁移切换。实现 WriteSyncer与 zapCore的接口。一次格式化，多次写入。
2. WriteSyncer接口实现，本地：日志的切割，压缩，保留期，动态切换等级等。远端：数据打包压缩传输，异常备份重试，传输方式支持 HTTP GRPC，网络带宽占用可控，内存占用低等特点。
3. 支持携带额外查询信息，比如 TraceId、多级查询条件筛选。
4. 系统的Receiver版块可横向扩展，保证高可用，高性能，高容量。
5. 报警模块，支持日志直接报警，可轻松制定关键词，报警频率，灵活开关，可实现多种报警方式(目前支持DingDing)。
6. 存储采用自维护分实例分库分集合功能。支持自定义天级别数据分片，项目分库。配置较为简单，支持库容量监控管理。可一直横向扩展，而不影响写入速度。每个项目存储日志周期可单独调节，自动管控容量，节约存储空间。特别优化联合索引查询，保证查询速度，降低索引大小。
7. 支持日志写入统计趋势，小时级别统计日志等级数量，主机写入数量等信息分析出一些有用信息。
8. 方便易操作的后台管理页面，直接访问(http://localhost:31080/admin)，让查询日志和配置管理都很高效。
9. 部署简单，依赖少，支持Docker快捷部署,开箱即用。

### 技术栈及设计简图

#### 技术栈:

1. 后端：语言Golang、协议GRPC HTTP、存储Mongodb
2. 前端：Vue

#### 设计简图

![设计简图](https://github.com/bbdshow/images/blob/master/qelog/qelog_design.png)

### 使用建议

#### 日志Client端

> go get -u github.com/bbdshow/qelog/qezap

配置文件参考 <a href="https://github.com/bbdshow/qelog/blob/main/configs/config.toml">configs/config.toml</a>

使用方式可以参考Qelog项目本身

<a href="https://github.com/bbdshow/qelog/blob/main/qezap/example/main.go">qezap/example/main.go</a>

#### 部署参考
- git clone 代码，进入项目根目录 
- 选择Docker部署，修改好 configs 里面的 config.docker.toml
```shell
# 构建镜像
make image
docker-compose up -d
```
- 选择nohup直接部署
```shell
# 构建镜像
make
cd ./bin
# 修改配置文件 configs/config.toml
# 可参照 /scripts start.sh stop.sh
# 也可以把 sh copy 到对应服务dir下 例如 admin/start.sh stop.sh
# 例如启动 admin
cd ./admin
nohup ./qelog_admin -f ../configs/config.toml >> nohup.out 2>&1  &
```
- 建议只启动一个 admin 即可，admin服务有定时任务

> 后台地址  http://localhost:31080/admin

#### 感谢支持,有问题请提 Issues，持续更新并解决问题。

