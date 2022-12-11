package types

import (
	"strings"
)

type Level int32

func (lvl Level) Int32() int32 {
	return int32(lvl)
}
func (lvl Level) String() string {
	v := "UNKNOWN"
	switch lvl {
	case -1:
		v = "DEBUG"
	case 0:
		v = "INFO"
	case 1:
		v = "WARN"
	case 2:
		v = "ERROR"
	case 3:
		v = "DPANIC"
	case 4:
		v = "PANIC"
	case 5:
		v = "FATAL"
	}
	return v
}

func String2Level(lvl string) Level {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return -1
	case "INFO":
		return 0
	case "WARN":
		return 1
	case "ERROR":
		return 2
	case "DPANIC":
		return 3
	case "PANIC":
		return 4
	case "FATAL":
		return 5
	}
	return -2
}
