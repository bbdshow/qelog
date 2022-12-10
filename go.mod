module github.com/bbdshow/qelog

go 1.14

replace (
	github.com/bbdshow/qelog/api => ./api
	github.com/bbdshow/qelog/qezap => ./qezap
)

require (
	github.com/bbdshow/bkit v0.3.7
	github.com/bbdshow/bkit/db/mongo v0.1.1
	github.com/bbdshow/qelog/api v1.1.0
	github.com/bbdshow/qelog/qezap v1.1.0
	github.com/gin-gonic/gin v1.7.2
	github.com/json-iterator/go v1.1.11
	go.mongodb.org/mongo-driver v1.7.2
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.41.0
)
