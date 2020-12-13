package receiver

import (
	"fmt"
	"time"

	"github.com/huzhongqing/qelog/model/mongoclient"

	"github.com/huzhongqing/qelog/storage"

	"github.com/huzhongqing/qelog/model"
	"github.com/huzhongqing/qelog/types/entity"
)

type Service struct {
	store *storage.Store
}

func NewService(database *mongoclient.Database) *Service {
	srv := &Service{store: storage.New(database)}

	return srv
}

func (srv *Service) InsertPacket(uk, ip string, packet entity.DataPacket) error {
	if len(packet.Data) <= 0 {
		return nil
	}
	records := srv.decodePacket(uk, ip, packet)
	fmt.Println(records)
	return nil
}

func (srv *Service) decodePacket(uk, ip string, packet entity.DataPacket) []*model.Logging {
	records := make([]*model.Logging, 0, len(packet.Data))
	for _, v := range packet.Data {
		r := &model.Logging{
			UniqueKey:     uk,
			Module:        packet.Name,
			IP:            ip,
			Full:          v,
			MillTimeStamp: time.Now().UnixNano() / 1e6,
		}
		val := make(map[string]interface{})
		if err := entity.Unmarshal([]byte(v), &val); err == nil {
			dec := entity.Decoder{Val: val}
			r.Short = dec.Short()
			r.Level = dec.Level()
			r.Condition1 = dec.Condition(1)
			r.Condition2 = dec.Condition(2)
			r.Condition3 = dec.Condition(3)
			r.MillTimeStamp = dec.Time()
		}
		records = append(records, r)
	}
	return records
}
