package mongo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// 分片集合名 生成规则
type ShardingCollection struct {
	prefix string
	// 天范围
	daySpan map[int]int
}

// NewShardingCollection
func NewShardingCollection(namePrefix string, span int) ShardingCollection {
	sn := ShardingCollection{prefix: namePrefix, daySpan: make(map[int]int)}
	sn.daySpan = sn.calcSpan(span)
	return sn
}

func (sc ShardingCollection) calcSpan(span int) map[int]int {
	// 31天，span=7 则 7天一个区间
	ds := span
	if ds <= 0 {
		ds = 1
	} else if ds >= 31 {
		ds = 31
	}
	day := 31
	size := 0
	s := 1
	daySpan := make(map[int]int, 0)
	for i := 1; i <= day; i++ {
		if size >= ds {
			size = 0
			s++
		}
		size++
		daySpan[i] = s
	}
	return daySpan
}

func (sc ShardingCollection) EncodeCollectionName(index int, unix int64) string {
	y, m, d := time.Unix(unix, 0).Date()
	s := sc.daySpan[d]
	return sc.collectionName(sc.prefix, index, y, int(m), s)
}

func (sc ShardingCollection) collectionName(prefix string, index, year, month, span int) string {
	v := fmt.Sprintf("%s_%d_%d%02d_%02d", prefix, index, year, month, span)
	return v
}

func (sc ShardingCollection) DecodeCollectionName(collectionName string) (prefix string, index, year int, month time.Month, span int, err error) {
	str := strings.Split(collectionName, "_")
	if len(str) != 4 {
		err = fmt.Errorf("invalid collection name %s", collectionName)
		return
	}
	prefix = str[0]
	di, _ := strconv.ParseInt(str[1], 10, 64)
	index = int(di)
	y, _ := strconv.ParseInt(str[2][:4], 10, 64)
	year = int(y)
	m, _ := strconv.ParseInt(str[2][4:], 10, 64)
	month = time.Month(m)
	s, _ := strconv.ParseInt(str[3], 10, 64)
	span = int(s)

	return
}

func (sc ShardingCollection) DaySpan() map[int]int {
	return sc.daySpan
}

func (sc ShardingCollection) CollectionNameToTime(collectionName string) (time.Time, error) {
	_, _, y, m, _, err := sc.DecodeCollectionName(collectionName)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(y, m, 0, 0, 0, 0, 0, time.Local), nil
}

func (sc ShardingCollection) SuggestSpanTime(collectionName string) (t time.Time, err error) {
	_, _, y1, m1, n1Span, err := sc.DecodeCollectionName(collectionName)
	if err != nil {
		return t, err
	}
	minDay := math.MaxInt32
	for d, span := range sc.daySpan {
		if span > n1Span {
			if d < minDay {
				minDay = d
			}
		}
	}
	t = time.Date(y1, m1, minDay, 0, 0, 0, 0, time.Local)
	return t, nil
}

// 根据开始时间和结束时间，查询出所有生成的 name
func (sc ShardingCollection) ScopeCollectionNames(index int, beginUnix, endUnix int64) []string {

	beginTime := time.Unix(beginUnix, 0)
	endTime := time.Unix(endUnix, 0)

	midTime := beginTime
	date := []time.Time{beginTime}
	for {
		midTime = midTime.AddDate(0, 0, 1)
		if midTime.Before(endTime) {
			date = append(date, midTime)
			continue
		}
		break
	}
	nameMap := make(map[string]struct{})
	names := make([]string, 0, len(date))
	for _, v := range date {
		name := sc.EncodeCollectionName(index, v.Unix())
		_, ok := nameMap[name]
		if !ok {
			nameMap[name] = struct{}{}
			names = append(names, name)
		}
	}
	return names
}
