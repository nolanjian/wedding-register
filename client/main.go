package main

import (
	"time"
	"wedding-register/request"

	"github.com/CodisLabs/codis/pkg/utils/log"
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
		log.Debug("fuck")
		return
	}

	// set now
	now := time.Now().In(local)

	// get deta
	tsub := timeBegin.Sub(now)

	log.Info(tsub.String(), " to begin")

	time.Sleep(tsub)

	log.Info("begin at: ", time.Now().In(local).Format(layout))
}

func main() {
	proxy := request.GetWebProxy()
	proxy.SetData("2017-05-19")

	timmer("04/19/2017 17:58:15 PM")

	for ii := 0; ii < 10; ii++ {
		err := proxy.Excute()
		if err != nil {
			log.Error(err)
		} else {
			log.Info(time.Now().In(local).Format(layout), " success")
			break
		}
	}
}
