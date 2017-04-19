package main

import (
	"time"
	"wedding-register/logger"
	"wedding-register/request"
)

var local *time.Location

// set layout
var layout = `01/02/2006 15:04:05 PM`

func timmer(beginTime string) {

	// set local
	local, _ = time.LoadLocation("Local")

	// set begin time
	timeBegin, err := time.ParseInLocation(layout, beginTime, local)
	if err != nil {
		logger.GetLogger().Debug("fuck")
		return
	}

	// set now
	now := time.Now().In(local)

	// get deta
	tsub := timeBegin.Sub(now)

	logger.GetLogger().Info(tsub.String(), " to begin")

	time.Sleep(tsub)

	logger.GetLogger().Info("begin at: ", time.Now().In(local).Format(layout))
}

func main() {
	proxy := request.GetWebProxy()
	proxy.SetData("2017-05-20")

	timmer("04/20/2017 00:00:00 AM")

	startStep := make([]int, 9)
	startStep[0] = 3
	startStep[1] = 0
	startStep[2] = 0
	startStep[3] = 0
	startStep[4] = 0
	startStep[5] = 0
	startStep[6] = 0
	startStep[7] = 0
	startStep[8] = 0

	const max = 10000

	sleepTime := make([]time.Duration, 9)
	sleepTime[0] = 18
	sleepTime[1] = 18
	sleepTime[2] = 9
	sleepTime[3] = 18
	sleepTime[4] = 28
	sleepTime[5] = 18
	sleepTime[6] = 18
	sleepTime[7] = 18
	sleepTime[8] = 9

	ii := 0
	for {
		logger.GetLogger().Info(time.Now().In(local).Format(layout), " start")

		ii = ii % max

		if err := proxy.Excute(startStep[ii%len(startStep)]); err == nil {
			logger.GetLogger().Info(time.Now().In(local).Format(layout), " success")
			break
		}

		logger.GetLogger().Info(time.Now().In(local).Format(layout), " sleep")

		time.Sleep(time.Second * sleepTime[ii%len(sleepTime)])

		ii++
	}
}
