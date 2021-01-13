module github.com/huzhongqing/qelog/qezap

go 1.14

replace github.com/huzhongqing/qelog/api => ../api

require (
	github.com/huzhongqing/qelog/api v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.34.1
)
