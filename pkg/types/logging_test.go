package types

import (
	"fmt"
	"testing"
	"time"
)

func TestLoggingCollectionName_FormatName(t *testing.T) {
	tStr := "20210221 08:00:00"
	now, err := time.Parse("20060102 15:04:05", tStr)
	if err != nil {
		t.Fatal(err)
	}
	n := NewLoggingCollectionName(7)
	name := n.FormatName(1, now.Unix())

	fmt.Println(name)
}

func TestLoggingCollectionName_ScopeNames(t *testing.T) {
	n := NewLoggingCollectionName(7)
	sStr := "20210321 08:00:00"
	start, err := time.Parse("20060102 15:04:05", sStr)
	if err != nil {
		t.Fatal(err)
	}
	eStr := "20210322 08:00:00"
	end, err := time.Parse("20060102 15:04:05", eStr)
	if err != nil {
		t.Fatal(err)
	}
	name := n.ScopeNames(1, start.Unix(), end.Unix())
	fmt.Println(name)
}

func TestLoggingCollectionName_SuggestTime(t *testing.T) {
	n := NewLoggingCollectionName(7)
	sStr := "20210322 08:00:00"
	start, err := time.Parse("20060102 15:04:05", sStr)
	if err != nil {
		t.Fatal(err)
	}
	eStr := "20210330 08:00:00"
	end, err := time.Parse("20060102 15:04:05", eStr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(n.DaySpan())
	name := n.ScopeNames(1, start.Unix(), end.Unix())
	fmt.Println(name)

	suggestTime, err := n.SuggestTime(name[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sStr, suggestTime.Add(-time.Second), "===", suggestTime, eStr)
}
