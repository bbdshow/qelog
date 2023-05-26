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
	mkdir -p ./bin/web && cp -r web ./bin/web
	go mod tidy
	@if [ ${OS} != "" ]; then\
		GOOS=${OS} go build -ldflags "-s" -o bin/qelog cmd/qelog/main.go;\
	else\
		go build -ldflags "-s" -o bin/qelog cmd/qelog/main.go;\
    fi

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./data
	rm -rf ./log
	rm -rf ./converage.txt

# build docker image
.PHONY: image
image:
	docker build -t ${tag} .

test:
	cd ./qezap && rm -r ./log||true && go mod tidy && go test -v -coverprofile=../converage.txt ./ && rm -r ./log
