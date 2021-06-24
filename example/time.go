package main

import (
	"fmt"
	"time"

	"github.com/haobird/goutils"
)

const NORMALTIMEFORMAT = "2006-01-02T15:04:05"

func main() {
	str := "2021-06-24T15:00:47"
	// temp, _ := goutils.GetTimeByNormalString(str)
	// timestamp := goutils.GetTimeUnix(temp)
	temp, _ := time.ParseInLocation(NORMALTIMEFORMAT, str, time.Local)
	timestamp := goutils.GetTimeUnix(temp)
	fmt.Println(timestamp)
}
