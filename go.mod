module github.com/huzhongqing/qelog

go 1.14

replace (
    github.com/huzhongqing/qelog/qezap => ./qezap
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.4.3
	github.com/huzhongqing/qelog/qezap v0.0.0-00010101000000-000000000000
	github.com/json-iterator/go v1.1.10
	go.mongodb.org/mongo-driver v1.4.4
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.34.1
)
