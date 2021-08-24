SET_OS=""
config=config.toml
tag=qelog:latest

.PHONY: build
build:
	export GOPROXY="https://goproxy.io,direct"
	mkdir -p ./bin && rm -r ./bin
	mkdir -p ./bin/configs && cp -r configs ./bin
	mkdir -p ./bin/admin/web && cp -r web ./bin/admin
	@if [ ${SET_OS} != "" ]; then\
		GOOS=${SET_OS} go build -o bin/receiver/qelog_receiver cmd/receiver/main.go;\
		GOOS=${SET_OS} go build -o bin/admin/qelog_admin cmd/admin/main.go;\
	else\
		go build -o bin/receiver/qelog_receiver cmd/receiver/main.go;\
		go build -o bin/admin/qelog_admin cmd/admin/main.go;\
    fi

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./data
	rm -rf ./log

# 构建Docker镜像
.PHONY: image
image:
	docker build -t ${tag} .
