# use GOOS build
OS=""
# default config file
config=config.toml
# default docker image tag version
tag=qelog:latest

.PHONY: build
build:
	export GOPROXY="https://goproxy.io,direct"
	mkdir -p ./bin && rm -r ./bin
	mkdir -p ./bin/configs && cp -r configs ./bin
	mkdir -p ./bin/admin/web && cp -r web ./bin/admin
	@if [ ${OS} != "" ]; then\
		GOOS=${OS} go build -ldflags "-s" -o bin/receiver/qelog_receiver cmd/receiver/main.go;\
		GOOS=${OS} go build -ldflags "-s" -o bin/admin/qelog_admin cmd/admin/main.go;\
	else\
		go build -ldflags "-s" -o bin/receiver/qelog_receiver cmd/receiver/main.go;\
		go build -ldflags "-s" -o bin/admin/qelog_admin cmd/admin/main.go;\
    fi

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./data
	rm -rf ./log

# build docker image
.PHONY: image
image:
	docker build -t ${tag} .
