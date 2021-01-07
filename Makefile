config=config.toml
tag=qelog:latest
flag="-s -w -X 'main.buildTime=`date`' -X 'main.goVersion=`go version`' -X 'main.gitHash=`git rev-parse HEAD`'"

.PHONY: build
build:
	mkdir -p ./bin && rm -r ./bin
	go build -ldflags ${flag} -o bin/receiver cmd/receiver/main.go
	go build -ldflags ${flag} -o bin/manager cmd/manager/main.go

.PHONY: clean
clean:
	rm -r ./bin

# 构建Docker镜像
.PHONY: buildImage
buildImage:
	docker build -t ${tag} .
