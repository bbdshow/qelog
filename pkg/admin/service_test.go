package admin

import (
	"github.com/bbdshow/qelog/tests"
)

func init() {
	tests.InitTestDepends("../../configs/config.docker.toml")
}

//func TestDropCollection(t *testing.T) {
//	svc := NewService()
//
//	err := svc.DropLoggingCollection(context.Background(), &entity.DropLoggingCollectionReq{
//		Host: "",
//		Name: "qelog_sharding2.logging_5_202101",
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//}
