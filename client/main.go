package main

import (
	"wedding/request"

	"github.com/CodisLabs/codis/pkg/utils/log"
)

func main() {
	proxy := request.GetWebProxy()
	proxy.SetData("2017-05-20")
	if err := proxy.Excute(); err != nil {
		log.Error(err)
	}
}
