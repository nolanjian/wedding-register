package main

import (
	"time"
	"wedding-register/request"

	"github.com/CodisLabs/codis/pkg/utils/log"
)

func timmer() {
	tm, _ := time.Parse("01/02/2006 15:04:05 PM", "05/20/2017 00:00:10 AM")
	begTime := tm.Unix()

	now := time.Now().Unix()

	sec := time.Unix(begTime-now, 0).Second()

	ss := time.Duration(sec) * time.Second

	log.Info(ss, " to begin")

	time.Sleep(ss)
	log.Info("begin at:", time.Now().Format("01/02/2006 15:04:05"))
}

func main() {
	proxy := request.GetWebProxy()
	proxy.SetData("2017-05-10")
	if err := proxy.Excute(); err != nil {
		log.Error(err)
	}
}
