config=config.toml
tag=qelog:latest
flag="-s -w -X 'main.buildTime=`date`' -X 'main.goVersion=`go version`' -X 'main.gitHash=`git rev-parse HEAD`'"

.PHONY: build
build:
	export GOPROXY="https://goproxy.io,direct"
	mkdir -p ./bin && rm -r ./bin
	mkdir -p ./bin/configs && cp -r configs ./bin
	mkdir -p ./bin/web && cp -r web ./bin
	go build -ldflags ${flag} -o bin/receiver cmd/receiver/main.go
	go build -ldflags ${flag} -o bin/manager cmd/manager/main.go

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./data
	rm -rf ./log

# 构建Docker镜像
.PHONY: buildImage
buildImage:
	docker build -t ${tag} .
