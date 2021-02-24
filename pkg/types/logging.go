package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type LoggingCollectionName struct {
	// 天范围
	daySpan map[int]int
}

func NewLoggingCollectionName(span int) LoggingCollectionName {
	lcn := LoggingCollectionName{daySpan: make(map[int]int)}
	lcn.daySpan = lcn.calcSpan(span)
	fmt.Println(lcn.daySpan)
	return lcn
}

func (lcn LoggingCollectionName) calcSpan(span int) map[int]int {
	ds := span
	if ds <= 0 {
		ds = 1
	} else if ds >= 31 {
		ds = 31
	}
	day := 31
	spanSize := day/ds + day%ds
	currentSpan := 1
	size := 0
	daySpan := make(map[int]int, 0)
	for i := 1; i <= day; i++ {
		size++
		daySpan[i] = currentSpan
		if size >= spanSize {
			currentSpan++
			size = 0
		}
	}
	return daySpan
}

func (lcn LoggingCollectionName) FormatName(index int, unix int64) string {
	y, m, d := time.Unix(unix, 0).Date()
	s := lcn.daySpan[d]
	return lcn.formatName("logging", index, y, int(m), s)
}

func (lcn LoggingCollectionName) formatName(prefix string, index, year, month, span int) string {
	v := fmt.Sprintf("%s_%d_%d%02d_%02d", prefix, index, year, month, span)
	return v
}

func (lcn LoggingCollectionName) Decode(name string) (prefix string, index, year, month, span int, err error) {
	str := strings.Split(name, "_")
	if len(str) != 4 {
		err = fmt.Errorf("invalid name %s", name)
		return
	}
	prefix = str[0]
	di, _ := strconv.ParseInt(str[1], 10, 64)
	index = int(di)
	y, _ := strconv.ParseInt(str[2][:4], 10, 64)
	year = int(y)
	m, _ := strconv.ParseInt(str[2][4:], 10, 64)
	month = int(m)
	s, _ := strconv.ParseInt(str[3], 10, 64)
	span = int(s)

	return
}

// 根据开始时间和结束时间，查询出所有生成的 name
func (lcn LoggingCollectionName) ScopeNames(index int, startUnix, endUnix int64) []string {
	// start end 生成所有 日期
	startTime := time.Unix(startUnix, 0)
	endTime := time.Unix(endUnix, 0)

	midTime := startTime
	date := []time.Time{startTime}
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
		name := lcn.FormatName(index, v.Unix())
		_, ok := nameMap[name]
		if !ok {
			nameMap[name] = struct{}{}
			names = append(names, name)
		}
	}
	return names
}
