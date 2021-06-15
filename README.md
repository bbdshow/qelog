# Qelog

Qelog是一款小巧且功能干练实用的日志系统(集成报警)。采用Push模式，拥有优秀的写入速度、可存储容量高，部署维护成本低等特性。其诞生的目的就是为了解决在中小团队或个人Go项目中，避免为日志系统投入过多的使用和运维成本(ELK)。

采用简单的设计做日志系统最本质的事情。已系统高效稳定运行、低使用成本及服务运维成本为目标。


### 日志系统特性:

1. 采用 Uber-zap 的Format，高效格式化日志，友好迁移切换。实现 WriteSyncer与 zapCore的接口。一次格式化，多次写入。
2. WriteSyncer接口实现，本地：日志的切割，压缩，保留期，动态切换等级等。远端：数据打包压缩传输，异常备份重试，传输方式支持 HTTP GRPC，网络带宽占用可控，内存占用低等特点。
3. 支持携带额外查询信息，比如 TraceId、多级查询条件筛选。
4. 系统的每个版块都可横向扩展，保证高可用，高性能，高容量。
5. 报警模块，支持日志直接报警，可轻松制定关键词，报警频率，灵活开关，可实现多种报警方式(目前DingDing)。
6. 存储采用自维护分实例分库分集合功能。支持自定义天级别数据分片，项目分库。配置较为简单，支持库容量监控管理。可一直横向扩展，而不影响写入速度。特别优化联合索引查询，保证查询速度，降低索引大小。
7. 支持日志写入统计趋势，小时级别统计日志等级数量，主机写入数量等信息分析出一些有用信息。
8. 方便快捷的前端管理平台（单页应用），让查询和配置日志更加高效。
9. 部署简单，依赖少，支持Dokcer快捷部署

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

<a href="https://github.com/bbdshow/qelog/blob/main/infra/logs/qezap.go">infra/logs/qezap.go</a> 

<a href="https://github.com/bbdshow/qelog/blob/main/infra/httputil/middleware.go">infra/httputil/middleware.go</a>  

<a href="https://github.com/bbdshow/qelog/blob/main/qezap/example/main.go">qezap/example/main.go</a>

### 后台部分截图

![查询](https://github.com/bbdshow/images/blob/master/qelog/find.png)
![报警](https://github.com/bbdshow/images/blob/master/qelog/alarm.png?raw=true)
![容量](https://github.com/bbdshow/images/blob/master/qelog/db.png)
![趋势](https://github.com/bbdshow/images/blob/master/qelog/trend.png)

**更多内容，还请部署后查看**

> 后台地址  http://localhost:31080/admin

#### 项目已线上稳定运行小半年，欢迎大家使用反馈问题提交PR等。
#### 感谢支持,喜欢点一个小星星，持续更新并解决问题。

