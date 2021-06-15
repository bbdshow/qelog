package admin

import (
	"context"
	"github.com/bbdshow/qelog/pkg/common/entity"
	"github.com/bbdshow/qelog/tests"
	"testing"
)

func init() {
	tests.InitTestDepends("../../configs/config.docker.toml")
}

func TestDropCollection(t *testing.T) {
	srv := NewService()

	err := srv.DropLoggingCollection(context.Background(), &entity.DropLoggingCollectionReq{
		Host: "",
		Name: "qelog_sharding2.logging_5_202101",
	})
	if err != nil {
		t.Fatal(err)
	}
}
