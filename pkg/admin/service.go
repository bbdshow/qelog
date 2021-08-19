package admin

import (
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/dao"

	"sync"
)

type Service struct {
	//sc                 mongo.ShardingCollection
	cfg  *conf.Config
	d    *dao.Dao
	once sync.Once
}

func NewService(cfg *conf.Config) *Service {
	svc := &Service{
		cfg: cfg,
		d:   dao.New(cfg),
		//sc:                 mongo.NewShardingCollection("logging", config.Global.DaySpan),
	}
	svc.once.Do(func() {
		//go svc.backgroundDelExpiredCollection(config.Global.MaxAgeMonth)
	})
	return svc
}

func (svc *Service) Close() {
	svc.d.Close()
}

//type AscShardingIndexState []entity.ShardingIndexState
//
//func (asc AscShardingIndexState) Len() int           { return len(asc) }
//func (asc AscShardingIndexState) Swap(i, j int)      { asc[i], asc[j] = asc[j], asc[i] }
//func (asc AscShardingIndexState) Less(i, j int) bool { return asc[i].Use < asc[j].Use }
//
//func (svc *Service) GetShardingIndex(ctx context.Context, out *entity.GetShardingIndexResp) error {
//	_, docs, err := svc.moduleStore.FindCountModule(ctx, bson.M{})
//	if err != nil {
//		return httputil.ErrSystemException.MergeError(err)
//	}
//	state := make(map[int]int)
//	for i := 1; i <= model.ShardingIndexSize; i++ {
//		state[i] = 0
//	}
//	for _, v := range docs {
//		num, ok := state[v.ShardingIndex]
//		if ok {
//			state[v.ShardingIndex] = num + 1
//		}
//	}
//
//	states := make([]entity.ShardingIndexState, 0, len(state))
//	for k, v := range state {
//		states = append(states, entity.ShardingIndexState{
//			Index: k,
//			Use:   v,
//		})
//	}
//	sort.Sort(AscShardingIndexState(states))
//
//	// 找到最小的，作为推荐
//	suggestDBIndex := model.ShardingIndexSize
//	if len(states) > 0 {
//		suggestDBIndex = states[0].Index
//	}
//
//	out.SuggestIndex = suggestDBIndex
//	out.ShardingIndexSize = model.ShardingIndexSize
//	out.UseState = states
//
//	return nil
//}
