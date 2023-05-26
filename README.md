# Qelog
[![Build&Test](https://github.com/bbdshow/qelog/actions/workflows/tests.yaml/badge.svg)](https://github.com/bbdshow/qelog/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbdshow/qelog)](https://goreportcard.com/report/github.com/bbdshow/qelog)
[![Go Reference](https://pkg.go.dev/badge/github.com/bbdshow/qelog/qezap.svg)](https://pkg.go.dev/github.com/bbdshow/qelog/qezap)
[![codecov](https://codecov.io/gh/bbdshow/qelog/branch/main/graph/badge.svg?token=Fqaz5qvx2Q)](https://codecov.io/gh/bbdshow/qelog)

Qelog is a small, low cost and light operation and maintenance of the log service. The purpose of its birth is to solve the problem of small and medium-sized team multi-service group log and alarm.

It has been running stably for 2+ years and has been rolling over 100+TB of data.

**[中文文档](./docs/README_CH.md)**

[admin manager example address](https://qelogdemo.bbdshow.top/admin)  username：admin passwd：111111

### Log System Features:

#### Client
- **qezap** is the Go version Client of the project. Wrap [Uber-zap](https://github.com/uber-go/zap) ,Extension for local and remote implementations of WriteSyncer.
- Local extension: log cutting, compression, retention period, dynamic switching level, etc.
- Remote expansion: The buffer is packed and compressed for transmission log data, and the error backup is retried to ensure that the data is not lost. 
The GRPC(default) HTTP protocol is supported, and the I/O control ensures the bandwidth usage is controllable and the memory usage is low.

#### Log receiver server
- The Receiver process can be expanded horizontally to ensure high availability and high performance.
- Configure multiple storage instances on the Receiver to improve the storage capacity and write performance of the cluster.
- Alarm module detects alarm rules for each log. The hit can be delivered according to the rules and different alarm methods. Currently supported DingTalk | Telegram
- Implement data fragmentation storage rules, support automatic capacity management, monitoring and early warning. Store separate instances of extensions without bottlenecks due to middleware.
- Log statistics, level distribution, and trend report.

#### Log manager server
- Based on the [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin) modification,support rich query dimensions, such as level, keyword, TraceId, ClientIP, multi-level conditions, etc. 
Due to cost and performance issues ** Full-text indexing ** is currently not supported
- Friendly operation interaction, simple configuration, efficient content display, almost can be **done out of the box quick used**.
- Cluster capacity statistics query, manual intervention, and other functions.

#### Performance test tool

- /tools Includes benchmark and memory footprint tests.
- Different operating environments and configurations have different performance. Therefore, you can conduct a comprehensive analysis to preliminarily check whether the service meets the requirements of scenarios.

#### Design drawing

![Design drawing](./docs/qelog_design.png)

### Usage

#### Qezap Client import your project

> go get -u github.com/bbdshow/qelog/qezap

[Client use example](./qezap/example/main.go)


#### Quick Deploy
default use single node deploy
##### Docker Deploy

```shell
git clone https://github.com/bbdshow/qelog.git
cd qelog
# build docker image
make image
# docker-compose start container
docker-compose up -d
```

##### Binary Deploy

```shell
git clone https://github.com/bbdshow/qelog.git
cd qelog
make
cd bin
# ./qelog or nohup ./qelog &
```

#### Service cluster deployment
Cluster deployment time, pay attention to mongo configuration file and qelog service running mode "cluster_admin" | "cluster_receiver", use docker-compose.yaml for cluster layout...

Thank you for your support. If it is useful to you, I hope the **Star** can support you. If you have any questions, please **Issues**, keep updating and solve problems.

