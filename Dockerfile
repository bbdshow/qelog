FROM golang:1.17.13-alpine3.16 as builder
WORKDIR /app

# when go.mod not change, use image cache
COPY go.mod .
COPY go.sum .

RUN set -eux; \
      ls -l; \
      export GOPROXY="https://goproxy.io,direct"; \
      go version && go mod download;

COPY . .

RUN set -eux; \
        mkdir -p ./bin && rm -r ./bin; \
        go mod tidy; \
        # shellcheck disable=SC2006
        go build -ldflags "-s" -o bin/qelog cmd/qelog/main.go;

FROM alpine:3.16
WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

## NOTE: tzdata default Asia/Shanghai
#RUN set -eux; \
#        apk update && apk add --no-cache\
#            tzdata \
#        ;\
#        cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime; \
#        echo 'Asia/Shanghai' > /etc/timezone; \
#        rm -rf /var/cache/apk/*;

COPY --from=builder ./app/bin .
COPY --from=builder ./app/configs ./configs
COPY --from=builder ./app/web ./web

ENTRYPOINT ["./qelog"]